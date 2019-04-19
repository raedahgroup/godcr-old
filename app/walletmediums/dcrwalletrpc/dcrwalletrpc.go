package dcrwalletrpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/decred/dcrd/wire"
	"github.com/decred/dcrwallet/netparams"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/dcrlibwallet/utils"
	"github.com/raedahgroup/godcr/app/config"
	"google.golang.org/grpc/codes"
)

// WalletRPCClient implements `WalletMiddleware` using `mobilewallet.LibWallet` as medium for connecting to a decred wallet
// Functions relating to operations that can be performed on a wallet are defined in `walletfunctions.go`
// Other wallet-related functions are defined in `walletloader.go`
type WalletRPCClient struct {
	walletLoader  walletrpc.WalletLoaderServiceClient
	walletService walletrpc.WalletServiceClient
	activeNet     *netparams.Params
	walletOpen    bool
}

// Connect establishes gRPC connection to a running dcrwallet daemon at the specified address,
// create a WalletServiceClient using the established connection. If the specified address did not connect,
// the RPC address is retreived from dcrwallet config file and if this fail, the default address is used.
// returns an instance of `dcrwalletrpc.Client`
func Connect(ctx context.Context, rpcAddress, rpcCert string, noTLS bool) (walletRPCClient *WalletRPCClient, err error) {
	defer func() {
		if walletRPCClient != nil {
			// wallet library is setup, prepare it for use by opening
			err = openWalletIfExist(ctx, walletRPCClient)
		}
	}()

	walletRPCClient, err = createConnection(ctx, rpcAddress, rpcCert, noTLS)
	if err == nil {
		return
	}

	walletRPCClient = parseDcrWalletConfigAndConnect(ctx, rpcCert, noTLS)
	if walletRPCClient != nil {
		return
	}

	walletRPCClient = connectToDefaultAddresses(ctx, rpcCert, noTLS)
	return
}

func createConnection(ctx context.Context, rpcAddress, rpcCert string, noTLS bool) (*WalletRPCClient, error) {
	if !noTLS && rpcCert == "" {
		return nil, errors.New("set dcrwallet rpc certificate path in config file or disable tls for dcrwallet connection")
	}

	// perform rpc connection in background, user might shutdown before connection is complete
	go connectToRPC(rpcAddress, rpcCert, noTLS)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	case connectionResult := <-rpcConnectionDone:
		if connectionResult.err != nil {
			return nil, connectionResult.err
		}

		return &WalletRPCClient{
			walletLoader:  walletrpc.NewWalletLoaderServiceClient(connectionResult.conn),
			walletService: walletrpc.NewWalletServiceClient(connectionResult.conn),
		}, nil
	}
}

func connectToDefaultAddresses(ctx context.Context, rpcCert string, noTLS bool) (walletRPCClient *WalletRPCClient) {
	// try connecting with default testnet3 params
	testnetAddress := net.JoinHostPort("localhost", netparams.TestNet3Params.GRPCServerPort)
	walletRPCClient, _ = createConnection(ctx, testnetAddress, rpcCert, noTLS)
	if walletRPCClient != nil {
		config.UpdateConfigFile(func(config *config.ConfFileOptions) {
			config.WalletRPCServer = testnetAddress
		})
		return
	}

	// try connecting with default mainnet params
	mainnetAddress := net.JoinHostPort("localhost", netparams.MainNetParams.GRPCServerPort)
	walletRPCClient, _ = createConnection(ctx, mainnetAddress, rpcCert, noTLS)
	if walletRPCClient != nil {
		config.UpdateConfigFile(func(config *config.ConfFileOptions) {
			config.WalletRPCServer = mainnetAddress
		})
		return
	}
	return
}

func openWalletIfExist(ctx context.Context, c *WalletRPCClient) error {
	c.walletOpen = false
	loadWalletDone := make(chan error)

	go func() {
		var openWalletError error
		defer func() {
			loadWalletDone <- openWalletError
		}()

		walletExists, openWalletError := c.WalletExists()
		if openWalletError != nil || !walletExists {
			return
		}

		_, openWalletError = c.walletLoader.OpenWallet(context.Background(), &walletrpc.OpenWalletRequest{})

		// ignore wallet already open errors, it could be that dcrwallet loaded the wallet when it was launched by the user
		// or godcr opened the wallet without closing it
		if isRpcErrorCode(openWalletError, codes.AlreadyExists) {
			openWalletError = nil
		}
	}()

	select {
	case err := <-loadWalletDone:
		// if err is nil, then wallet was opened
		if err == nil {
			c.walletOpen = true
			// wallet is open, best time to detect network type for dcrwallet rpc connection
			c.activeNet, _ = getNetParam(c.walletService)
		} else {
			c.walletOpen = false
		}
		return err

	case <-ctx.Done():
		return ctx.Err()
	}
}

func getNetParam(walletService walletrpc.WalletServiceClient) (param *netparams.Params, err error) {
	req := &walletrpc.NetworkRequest{}
	res, err := walletService.Network(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("error checking wallet rpc network type: %s", err.Error())
	}

	param = utils.NetParams(wire.CurrencyNet(res.ActiveNetwork).String())
	if param == nil {
		err = fmt.Errorf("unknown network type")
	}
	return
}

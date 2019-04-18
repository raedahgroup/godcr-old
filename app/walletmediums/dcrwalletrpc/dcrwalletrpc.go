package dcrwalletrpc

import (
	"context"
	"errors"
	"net"

	"github.com/decred/dcrwallet/netparams"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/godcr/app/config"
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
func Connect(ctx context.Context, rpcAddress, rpcCert string, noTLS bool) (*WalletRPCClient, error) {
	walletRPCClient, originalConnectionError := createConnection(ctx, rpcAddress, rpcCert, noTLS)
	if originalConnectionError == nil {
		return walletRPCClient, nil
	}

	walletRPCClient = parseDcrWalletConfigAndConnect(ctx, rpcCert, noTLS)
	if walletRPCClient != nil {
		return walletRPCClient, nil
	}

	walletRPCClient = connectToDefaultAddresses(ctx, rpcCert, noTLS)
	if walletRPCClient != nil {
		return walletRPCClient, nil
	}

	return nil, originalConnectionError
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

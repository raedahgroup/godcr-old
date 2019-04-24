package dcrwalletrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/decred/dcrd/wire"
	"github.com/decred/dcrwallet/netparams"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
	"github.com/raedahgroup/dcrlibwallet/txindex"
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
	walletOpen    bool
	activeNet     *netparams.Params

	numberOfPeers int32
	syncListener  *defaultsynclistener.DefaultSyncListener

	txIndexDB              *txindex.DB
	txNotificationListener TransactionListener
}

// Connect establishes gRPC connection to a running dcrwallet daemon at the specified address,
// create a WalletServiceClient using the established connection. If the specified address did not connect,
// the RPC address is retreived from dcrwallet config file and if this fail, the default address is used.
// returns an instance of `dcrwalletrpc.Client`
func Connect(ctx context.Context, cfg *config.Config) (walletRPCClient *WalletRPCClient, err error) {
	defer func() {
		if walletRPCClient != nil {
			// wallet library is setup, prepare it for use by opening
			err = openWalletIfExist(ctx, walletRPCClient, cfg.AppDataDir)
		}
	}()

	walletRPCClient, err = createConnection(ctx, cfg.WalletRPCServer, cfg.WalletRPCCert, cfg.NoWalletRPCTLS)
	if err == nil {
		return
	}

	walletRPCClient = parseDcrWalletConfigAndConnect(ctx, cfg.WalletRPCCert, cfg.NoWalletRPCTLS)
	if walletRPCClient != nil {
		return
	}

	walletRPCClient = connectToDefaultAddresses(ctx, cfg.WalletRPCCert, cfg.NoWalletRPCTLS)
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

func openWalletIfExist(ctx context.Context, c *WalletRPCClient, appDataDir string) error {
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
			err = finalizeWalletSetup(ctx, c, appDataDir)
		} else {
			c.walletOpen = false
		}
		return err

	case <-ctx.Done():
		return ctx.Err()
	}
}

func finalizeWalletSetup(ctx context.Context, c *WalletRPCClient, appDataDir string) error {
	c.walletOpen = true

	// wallet is open, best time to detect network type for dcrwallet rpc connection
	c.activeNet, _ = getNetParam(c.walletService)

	// set database for indexing transactions for faster loading
	// important to do it at this point before wallet operations
	// such as sync and transaction notification are triggered
	// because those operations will need to access the tx index db.
	txIndexDbPath := filepath.Join(appDataDir, "rpc-tx-index", txindex.DbName)
	os.MkdirAll(filepath.Dir(txIndexDbPath), os.ModePerm) // create directory if not exist

	generateWalletAddress := func() (string, error) {
		return c.GenerateNewAddress(0) // use default account
	}
	addressMatchesWallet := func(address string) (bool, error) {
		addressInfo, err := c.AddressInfo(address)
		if err != nil {
			return false, err
		}
		return addressInfo.IsMine, nil
	}

	txIndexDB, err := txindex.Initialize(txIndexDbPath, generateWalletAddress, addressMatchesWallet)
	if err != nil {
		return fmt.Errorf("tx index db initialization failed: %s", err.Error())
	}
	c.txIndexDB = txIndexDB

	// start tx notification listener now,
	// so we can index txs as the wallet is notified of new/updated txs
	c.ListenForTxNotification(ctx)

	return nil
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

package dcrwalletrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/decred/dcrwallet/netparams"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/godcr/app/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

type rpcConnectionResult struct {
	err  error
	conn *grpc.ClientConn
}

var (
	rpcConnectionDone    = make(chan *rpcConnectionResult)
	rpcConnectionTimeout = 5 * time.Second
)

// New establishes gRPC connection to a running dcrwallet daemon at the specified address,
// create a WalletServiceClient using the established connection. If the specified address did not connect,
// the RPC address is retreived from dcrwallet config file and if this fail, the default address is used.
// returns an instance of `dcrwalletrpc.Client`
func New(ctx context.Context, rpcAddress, rpcCert string, noTLS bool) (*WalletRPCClient, error) {
	walletRPCClient, originalConnectionError := createConnection(ctx, rpcAddress, rpcCert, noTLS)
	if originalConnectionError == nil {
		return walletRPCClient, originalConnectionError
	}

	dcrwalletConfAddresses, dcrwalletConfNoTLS, dcrwalletConfCert, err := connectionParamsFromDcrwalletConfig()
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to %s. Trying to check dcrwallet config for different address failed. %s", rpcAddress, err.Error())
	}

	if dcrwalletConfCert != "" {
		rpcCert = dcrwalletConfCert
	}
	noTLS = dcrwalletConfNoTLS

	if walletRPCClient = useParsedConfigAddresses(ctx, dcrwalletConfAddresses, rpcCert, noTLS); walletRPCClient != nil {
		return walletRPCClient, nil
	}

	if walletRPCClient = connectToDefaultAddresses(ctx, rpcCert, noTLS); walletRPCClient != nil {
		return walletRPCClient, nil
	}

	return nil, originalConnectionError
}

func useParsedConfigAddresses(ctx context.Context, addresses []string, rpcCert string, noTLS bool) (walletRPCClient *WalletRPCClient) {
	for _, address := range addresses {
		walletRPCClient, _ = createConnection(ctx, address, rpcCert, noTLS)
		if walletRPCClient != nil {
			config.UpdateConfigFile(func(config *config.ConfFileOptions) {
				config.WalletRPCServer = address
			})
			return
		}
	}
	return
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

		walletService := walletrpc.NewWalletServiceClient(connectionResult.conn)

		client := &WalletRPCClient{
			walletLoader:  walletrpc.NewWalletLoaderServiceClient(connectionResult.conn),
			walletService: walletService,
		}

		return client, nil
	}
}

func connectToRPC(rpcAddress, rpcCert string, noTLS bool) {
	var conn *grpc.ClientConn
	var err error

	defer func() {
		if conn == nil && err == nil {
			// connection timeout
			err = fmt.Errorf("Error connecting to %s. Connection attempt timed out after %s", rpcAddress, rpcConnectionTimeout)
		}
		connectionResult := &rpcConnectionResult{
			err:  err,
			conn: conn,
		}
		rpcConnectionDone <- connectionResult
	}()

	// block until connection is established
	// return error if connection cannot be established after `rpcConnectionTimeoutSeconds` seconds
	grpcConnectionOptions := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(rpcConnectionTimeout),
	}

	if noTLS {
		grpcConnectionOptions = append(grpcConnectionOptions, grpc.WithInsecure())
		conn, err = grpc.Dial(rpcAddress, grpcConnectionOptions...)
	} else {
		creds, err := credentials.NewClientTLSFromFile(rpcCert, "")
		if err != nil {
			return
		}

		grpcConnectionOptions = append(grpcConnectionOptions, grpc.WithTransportCredentials(creds))
		conn, err = grpc.Dial(rpcAddress, grpcConnectionOptions...)
	}
}

package dcrwalletrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/wire"
	"github.com/decred/dcrdata/netparams"
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
	activeNet     *chaincfg.Params
	walletOpen    bool
}

type rpcConnectionResult struct {
	err  error
	conn *grpc.ClientConn
}

const dcrwalletExecutableName = "dcrwallet"

var (
	rpcConnectionDone    = make(chan *rpcConnectionResult)
	rpcConnectionTimeout = 5 * time.Second
	errSkip              = errors.New("skip")
)

// New establishes gRPC connection to a running dcrwallet daemon at the specified address,
// create a WalletServiceClient using the established connection and
// returns an instance of `dcrwalletrpc.Client`
func New(ctx context.Context, rpcAddress, rpcCert string, noTLS bool) (*WalletRPCClient, error) {
	client, err1 := createConnection(ctx, rpcAddress, rpcCert, noTLS)
	if err1 == nil {
		return client, err1
	}
	rpcConfAddresses, wnoTLS, wrpcCert, err := walletAddressFromDcrdwalletConfig()
	if err != nil {
		return nil, err
	}

	if wrpcCert != "" {
		rpcCert = wrpcCert
	}
	noTLS = wnoTLS

	if client, err = useParsedConfigAddresses(rpcConfAddresses, ctx, rpcCert, noTLS); err != errSkip {
		return client, err
	}

	if client, err = connectToDefaultAddresses(ctx, rpcCert, noTLS); err != errSkip {
		return client, err
	}

	return nil, err1
}

func useParsedConfigAddresses(addresses []string, ctx context.Context, rpcCert string, noTLS bool) (client *WalletRPCClient, err error) {
	for _, address := range addresses {
		walletMiddleware, err := createConnection(ctx, address, rpcCert, noTLS)
		if err == nil {
			saveRPCAddress(address)
			return walletMiddleware, nil
		}
	}
	return nil, errSkip
}

func connectToDefaultAddresses(ctx context.Context, rpcCert string, noTLS bool) (*WalletRPCClient, error) {
	// try connecting with default testnet3 params
	testnetAddress := net.JoinHostPort("localhost", netparams.TestNet3Params.GRPCServerPort)
	walletMiddleware, err := createConnection(ctx, testnetAddress, rpcCert, noTLS)
	if err == nil {
		saveRPCAddress(testnetAddress)
		return walletMiddleware, nil
	}

	// try connecting with default mainnet params
	mainnetAddress := net.JoinHostPort("localhost", netparams.MainNetParams.GRPCServerPort)
	walletMiddleware, err = createConnection(ctx, mainnetAddress, rpcCert, noTLS)
	if err == nil {
		saveRPCAddress(mainnetAddress)
		return walletMiddleware, nil
	}
	return nil, errSkip
}

func createConnection(ctx context.Context, rpcAddress, rpcCert string, noTLS bool) (*WalletRPCClient, error) {
	// check if user has provided enough information to attempt connecting to dcrwallet
	if rpcAddress == "" {
		return nil, errors.New("you must set walletrpcserver in config file to use wallet rpc")
	}
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

		activeNet, err := getNetParam(walletService)
		if err != nil {
			return nil, err
		}

		client := &WalletRPCClient{
			walletLoader:  walletrpc.NewWalletLoaderServiceClient(connectionResult.conn),
			walletService: walletService,
			activeNet:     activeNet,
		}

		return client, nil
	}
}

func saveRPCAddress(address string) {
	err := config.UpdateConfigFile("walletrpcserver", address, true)
	if err != nil {
		fmt.Printf("unable to save rpc address: %s", err.Error())
	}
}

func getNetParam(walletService walletrpc.WalletServiceClient) (param *chaincfg.Params, err error) {
	req := &walletrpc.NetworkRequest{}
	res, err := walletService.Network(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("error checking wallet rpc network type: %s", err.Error())
	}

	switch res.GetActiveNetwork() {
	case uint32(wire.MainNet):
		return &chaincfg.MainNetParams, nil
	case uint32(wire.TestNet3):
		return &chaincfg.TestNet3Params, nil
	default:
		return nil, errors.New("unknown network type")
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

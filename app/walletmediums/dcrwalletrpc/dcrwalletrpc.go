package dcrwalletrpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ademuanthony/ps"
	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/wire"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
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

const (
	dcrwalletExecutableName = "dcrwallet"
)

var (
	rpcConnectionDone    = make(chan *rpcConnectionResult)
	rpcConnectionTimeout = 5 * time.Second
)

// New establishes gRPC connection to a running dcrwallet daemon at the specified address,
// create a WalletServiceClient using the established connection and
// returns an instance of `dcrwalletrpc.Client`
func New(ctx context.Context, rpcAddress, rpcCert string, noTLS bool) (client *WalletRPCClient, err error) {
	client, err = createConnection(ctx, rpcAddress, rpcCert, noTLS)
	if err == nil || dcrwalletIsRunning() {
		return
	}

	prompt := "There is no dcrwallet process started, would you like godcr to run the process"
	startWalletConfirmed, err := terminalprompt.RequestYesNoConfirmation(prompt, "Y")
	if err == nil {
		return nil, fmt.Errorf("error in reading input %s", err.Error())
	}
	if !startWalletConfirmed {
		return
	}
	err = findAndStartDcrwallet()
	if err != nil {
		return nil, err
	}
	return createConnection(ctx, rpcAddress, rpcCert, noTLS)
}

func findAndStartDcrwallet() error {
	panic("not implemented")
}

func dcrwalletIsRunning() bool {
	proc, _ := ps.ProcessByName(dcrwalletExecutableName)
	if proc == nil {
		return false
	}
	return true
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

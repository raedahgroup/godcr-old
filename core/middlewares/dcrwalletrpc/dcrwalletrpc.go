package dcrwalletrpc

import (
	"fmt"
	"github.com/decred/dcrwallet/netparams"
	pb "github.com/decred/dcrwallet/rpc/walletrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

// WalletPRCClient implements `WalletSource` using dcrwallet's `walletrpc.WalletServiceClient`
// Method implementation of `WalletSource` interface are in functions.go
// Other functions not related to `WalletSource` are in utils.go
type WalletPRCClient struct {
	walletService pb.WalletServiceClient
	activeNet     *chaincfg.Params
}

// New establishes gRPC connection to a running dcrwallet daemon at the specified address,
// create a WalletServiceClient using the established connection and
// returns an instance of `dcrwalletrpc.Client`
func New(address, cert string, noTLS, isTestnet bool) (*WalletPRCClient, error) {
	if rpcAddress == "" {
		rpcAddress = defaultDcrWalletRPCAddress(netType)
		fmt.Println(rpcAddress)
	}

	conn, err := connectToRPC(rpcAddress, rpcCert, noTLS)
	if err != nil {
		return nil, err
	}

	activeNet := &chaincfg.MainNetParams
	if isTestnet {
		activeNet = &chaincfg.TestNet3Params
	}

	client := &WalletPRCClient{
		walletService: pb.NewWalletServiceClient(conn),
		activeNet: activeNet,
	}

	return client, nil
}

func defaultDcrWalletRPCAddress(netType string) string {
	var activeNet *netparams.Params
	if netType == "mainnet" {
		activeNet = &netparams.MainNetParams
	} else {
		activeNet = &netparams.TestNet3Params
	}

	return net.JoinHostPort("localhost", activeNet.GRPCServerPort)
}

// todo remember to close grpc connection after usage
func connectToRPC(rpcAddress, rpcCert string, noTLS bool) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error

	if noTLS {
		conn, err = grpc.Dial(rpcAddress, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
	} else {
		creds, err := credentials.NewClientTLSFromFile(rpcCert, "")
		if err != nil {
			return nil, err
		}

		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(
				creds,
			),
		}

		conn, err = grpc.Dial(rpcAddress, opts...)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

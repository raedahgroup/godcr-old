package dcrwalletrpc

import (
	"net"

	"github.com/decred/dcrwallet/netparams"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// WalletPRCClient implements `WalletMiddleware` using `mobilewallet.LibWallet` as medium for connecting to a decred wallet
// Functions relating to operations that can be performed on a wallet are defined in `walletfunctions.go`
// Other wallet-related functions are defined in `walletloader.go`
type WalletPRCClient struct {
	walletLoader walletrpc.WalletLoaderServiceClient
	walletService walletrpc.WalletServiceClient
	netType       string
}

func New(netType, rpcAddress, rpcCert string, noTLS bool) (*WalletPRCClient, error) {
	if rpcAddress == "" {
		rpcAddress = defaultDcrWalletRPCAddress(netType)
	}

	conn, err := connectToRPC(rpcAddress, rpcCert, noTLS)
	if err != nil {
		return nil, err
	}

	client := &WalletPRCClient{
		walletLoader:  walletrpc.NewWalletLoaderServiceClient(conn),
		walletService: walletrpc.NewWalletServiceClient(conn),
		netType:       netType,
	}

	return client, nil
}

func defaultDcrWalletRPCAddress(netType string) string {
	if netType == "mainnet" {
		return net.JoinHostPort("localhost", netparams.MainNetParams.GRPCServerPort)
	} else {
		return net.JoinHostPort("localhost", netparams.TestNet3Params.GRPCServerPort)
	}
}

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

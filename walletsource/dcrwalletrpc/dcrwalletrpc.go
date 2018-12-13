package dcrwalletrpc

import (
	pb "github.com/decred/dcrwallet/rpc/walletrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// WalletPRCClient implements `WalletSource` using dcrwallet's `walletrpc.WalletServiceClient`
// Method implementation of `WalletSource` interface are in functions.go
// Other functions not related to `WalletSource` are in helpers.go
type WalletPRCClient struct {
	walletService pb.WalletServiceClient
}

func New(address, cert string, noTLS bool) (*WalletPRCClient, error) {
	conn, err := connectToRPC(address, cert, noTLS)
	if err != nil {
		return nil, err
	}

	client := &WalletPRCClient{
		walletService: pb.NewWalletServiceClient(conn),
	}

	return client, nil
}

// todo remember to close grpc connection after usage
func connectToRPC(address, cert string, noTLS bool) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error

	if noTLS {
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
	} else {
		creds, err := credentials.NewClientTLSFromFile(cert, "")
		if err != nil {
			return nil, err
		}

		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(
				creds,
			),
		}

		conn, err = grpc.Dial(address, opts...)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

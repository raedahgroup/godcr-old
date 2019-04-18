package dcrwalletrpc

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type rpcConnectionResult struct {
	err  error
	conn *grpc.ClientConn
}

var (
	rpcConnectionDone    = make(chan *rpcConnectionResult)
	rpcConnectionTimeout = 5 * time.Second
)

func connectToRPC(rpcAddress, rpcCert string, noTLS bool) {
	var conn *grpc.ClientConn
	var err error

	defer func() {
		if conn == nil && err == nil {
			// connection timeout
			err = fmt.Errorf("cannot connect to the rpc address '%s', connection attempt timed out after %s", rpcAddress,
				rpcConnectionTimeout)
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

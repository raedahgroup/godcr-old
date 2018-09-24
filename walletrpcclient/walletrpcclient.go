package walletrpcclient

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type (
	Response struct {
		Columns []string
		Result  []string
	}
	Handler func(conn *grpc.ClientConn, ctx context.Context, args []string) (*Response, error)
)

var (
	funcMap  = map[string]Handler{}
	commands = map[string]string{} // map of supported commands and description
)

// Connect attempts connection toi a gRPC server using the passed options
func Connect(address, cert string, noTLS bool) (*grpc.ClientConn, error) {
	if noTLS {
		return grpc.Dial(address, grpc.WithInsecure())
	}

	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		return nil, err
	}

	// dial options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(
			creds,
		),
	}
	return grpc.Dial(address, opts...)
}

// IsCommandSupported returns a boolean whose value depends on if a command is registered as suppurted along
// with it's func handler
func IsCommandSupported(command string) bool {
	_, ok := funcMap[command]
	return ok
}

// RunCommand takes a command and tries to call the appropriate handler to call a gRPC service
// This should only be called after verifying that the command is supported using the IsCommandSupported
// function.
func RunCommand(conn *grpc.ClientConn, command string, opts []string) (*Response, error) {
	handler := funcMap[command]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	res, err := handler(conn, ctx, opts)
	return res, err
}

// RegisterHandler registers a command, its description and its handler
func RegisterHandler(key, description string, h Handler) {
	if _, ok := funcMap[key]; ok {
		panic("trying to register a handler twice: " + key)
	}

	funcMap[key] = h
	commands[key] = description
}

package dcrwalletrpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func isRpcErrorCode(err error, code codes.Code) bool {
	if err == nil {
		return false
	}

	e, ok := status.FromError(err)
	return ok && e.Code() == code
}
package dcrwalletrpc

import (
	"testing"

	"google.golang.org/grpc/codes"
)

func Test_isRpcErrorCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code codes.Code
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRpcErrorCode(tt.err, tt.code); got != tt.want {
				t.Errorf("isRpcErrorCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

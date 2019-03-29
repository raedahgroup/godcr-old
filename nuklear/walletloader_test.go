package nuklear

import (
	"context"
	"testing"

	"github.com/raedahgroup/godcr/app"
)

func Test_openWalletIfExist(t *Testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		walletMiddleware app.WalletMiddleware
		wantWalletExists bool
		wantErr          bool
	}{
		// TODO: add test cases
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotWalletExists, err := openWalletIfExist(test.ctx, test.walletMiddleware)
			if (err != nil) != test.wantErr {
				t.Errorf("openWalletIfExist() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if gotWalletExists != test.wantWalletExists {
				t.Errorf("openWalletIfExist() = %v, want %v", gotWalletExists, test.wantWalletExists)
			}
		})
	}
}

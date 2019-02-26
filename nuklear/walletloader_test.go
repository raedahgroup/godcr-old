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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWalletExists, err := openWalletIfExist(tt.ctx, tt.walletMiddleware)
			if (err != nil) != tt.wantErr {
				t.Errorf("openWalletIfExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWalletExists != tt.wantWalletExists {
				t.Errorf("openWalletIfExist() = %v, want %v", gotWalletExists, tt.wantWalletExists)
			}
		})
	}
}

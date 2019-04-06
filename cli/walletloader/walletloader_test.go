package walletloader

import (
	"context"
	"testing"

	"github.com/raedahgroup/godcr/app"
)

func TestOpenWallet(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		walletMiddleware app.WalletMiddleware
		wantWalletExists bool
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotWalletExists, err := OpenWallet(test.ctx, test.walletMiddleware)
			if (err != nil) != test.wantErr {
				t.Errorf("OpenWallet() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if gotWalletExists != test.wantWalletExists {
				t.Errorf("OpenWallet() = %v, want %v", gotWalletExists, test.wantWalletExists)
			}
		})
	}
}

func TestSyncBlockChain(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		walletMiddleware app.WalletMiddleware
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := SyncBlockChain(test.ctx, test.walletMiddleware); (err != nil) != test.wantErr {
				t.Errorf("SyncBlockChain() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

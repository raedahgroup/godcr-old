package terminal

import (
	"context"
	"testing"

	"github.com/raedahgroup/godcr/app"
)

func Test_openWalletIfExist(t *testing.T) {
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

func TestCreateWallet(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		password         string
		walletMiddleware app.WalletMiddleware
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := CreateWallet(test.ctx, test.password, test.walletMiddleware); (err != nil) != test.wantErr {
				t.Errorf("CreateWallet() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

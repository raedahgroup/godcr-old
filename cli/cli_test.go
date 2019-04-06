package cli

import (
	"context"
	"testing"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		walletMiddleware app.WalletMiddleware
		appConfig        *config.Config
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := Run(test.ctx, test.walletMiddleware, test.appConfig); (err != nil) != test.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func Test_syncBlockChain(t *testing.T) {
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
			if err := syncBlockChain(test.ctx, test.walletMiddleware); (err != nil) != test.wantErr {
				t.Errorf("syncBlockChain() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func Test_listCommands(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			listCommands()
		})
	}
}

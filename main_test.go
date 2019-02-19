package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
)

func Test_attemptExecuteSimpleOp(t *testing.T) {
	tests := []struct {
		name           string
		wantIsSimpleOp bool
		wantErr        bool
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsSimpleOp, err := attemptExecuteSimpleOp()
			if (err != nil) != tt.wantErr {
				t.Errorf("attemptExecuteSimpleOp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsSimpleOp != tt.wantIsSimpleOp {
				t.Errorf("attemptExecuteSimpleOp() = %v, want %v", gotIsSimpleOp, tt.wantIsSimpleOp)
			}
		})
	}
}

func Test_connectToWallet(t *testing.T) {
	tests := []struct {
		name                 string
		ctx                  context.Context
		config               config.Config
		wantWalletMiddleware app.WalletMiddleware
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWalletMiddleware := connectToWallet(tt.ctx, tt.config); !reflect.DeepEqual(gotWalletMiddleware, tt.wantWalletMiddleware) {
				t.Errorf("connectToWallet() = %v, want %v", gotWalletMiddleware, tt.wantWalletMiddleware)
			}
		})
	}
}

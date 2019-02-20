package web

import (
	"context"
	"testing"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/godcr/app"
)

func TestStartServer(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		walletMiddleware app.WalletMiddleware
		host             string
		port             string
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StartServer(tt.ctx, tt.walletMiddleware, tt.host, tt.port); (err != nil) != tt.wantErr {
				t.Errorf("StartServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_openWalletIfExists(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		walletMiddleware app.WalletMiddleware
		wantErr          bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := openWalletIfExist(tt.ctx, tt.walletMiddleware); (err != nil) != tt.wantErr {
				t.Errorf("openWalletIfExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_startServer(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		address string
		router  chi.Router
		wantErr bool
	}{
		// TODO: add test cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := startServer(tt.ctx, tt.address, tt.router); (err != nil) != tt.wantErr {
				t.Errorf("startServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

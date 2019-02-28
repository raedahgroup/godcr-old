package dcrwalletrpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrwallet/rpc/walletrpc"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		rpcAddress string
		rpcCert    string
		noTLS      bool
		want       *WalletRPCClient
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.ctx, tt.rpcAddress, tt.rpcCert, tt.noTLS)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_useParsedConfigAddresses(t *testing.T) {
	tests := []struct {
		name                string
		ctx                 context.Context
		addresses           []string
		rpcCert             string
		noTLS               bool
		wantWalletRPCClient *WalletRPCClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWalletRPCClient := useParsedConfigAddresses(tt.ctx, tt.addresses, tt.rpcCert, tt.noTLS); !reflect.DeepEqual(gotWalletRPCClient, tt.wantWalletRPCClient) {
				t.Errorf("useParsedConfigAddresses() = %v, want %v", gotWalletRPCClient, tt.wantWalletRPCClient)
			}
		})
	}
}

func Test_connectToDefaultAddresses(t *testing.T) {
	tests := []struct {
		name                string
		ctx                 context.Context
		rpcCert             string
		noTLS               bool
		wantWalletRPCClient *WalletRPCClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWalletRPCClient := connectToDefaultAddresses(tt.ctx, tt.rpcCert, tt.noTLS); !reflect.DeepEqual(gotWalletRPCClient, tt.wantWalletRPCClient) {
				t.Errorf("connectToDefaultAddresses() = %v, want %v", gotWalletRPCClient, tt.wantWalletRPCClient)
			}
		})
	}
}

func Test_createConnection(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		rpcAddress string
		rpcCert    string
		noTLS      bool
		want       *WalletRPCClient
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createConnection(tt.ctx, tt.rpcAddress, tt.rpcCert, tt.noTLS)
			if (err != nil) != tt.wantErr {
				t.Errorf("createConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_connectToRPC(t *testing.T) {
	tests := []struct {
		name       string
		rpcAddress string
		rpcCert    string
		noTLS      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connectToRPC(tt.rpcAddress, tt.rpcCert, tt.noTLS)
		})
	}
}

func Test_getNetParam(t *testing.T) {
	tests := []struct {
		name          string
		walletService walletrpc.WalletServiceClient
		wantParam     *chaincfg.Params
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParam, err := getNetParam(tt.walletService)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNetParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotParam, tt.wantParam) {
				t.Errorf("getNetParam() = %v, want %v", gotParam, tt.wantParam)
			}
		})
	}
}

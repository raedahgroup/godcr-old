package dcrwalletrpc

import (
	"testing"

	"github.com/decred/dcrwallet/netparams"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/godcr/app"
)

func TestWalletRPCClient_NetType(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &WalletRPCClient{
				walletLoader:  tt.fields.walletLoader,
				walletService: tt.fields.walletService,
				activeNet:     tt.fields.activeNet,
				walletOpen:    tt.fields.walletOpen,
			}
			if got := c.NetType(); got != tt.want {
				t.Errorf("WalletRPCClient.NetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletRPCClient_WalletExists(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &WalletRPCClient{
				walletLoader:  tt.fields.walletLoader,
				walletService: tt.fields.walletService,
				activeNet:     tt.fields.activeNet,
				walletOpen:    tt.fields.walletOpen,
			}
			got, err := c.WalletExists()
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletRPCClient.WalletExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WalletRPCClient.WalletExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletRPCClient_GenerateNewWalletSeed(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &WalletRPCClient{
				walletLoader:  tt.fields.walletLoader,
				walletService: tt.fields.walletService,
				activeNet:     tt.fields.activeNet,
				walletOpen:    tt.fields.walletOpen,
			}
			got, err := c.GenerateNewWalletSeed()
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletRPCClient.GenerateNewWalletSeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WalletRPCClient.GenerateNewWalletSeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletRPCClient_CreateWallet(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
	}
	tests := []struct {
		name       string
		fields     fields
		passphrase string
		seed       string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &WalletRPCClient{
				walletLoader:  tt.fields.walletLoader,
				walletService: tt.fields.walletService,
				activeNet:     tt.fields.activeNet,
				walletOpen:    tt.fields.walletOpen,
			}
			if err := c.CreateWallet(tt.passphrase, tt.seed); (err != nil) != tt.wantErr {
				t.Errorf("WalletRPCClient.CreateWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWalletRPCClient_OpenWallet(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &WalletRPCClient{
				walletLoader:  tt.fields.walletLoader,
				walletService: tt.fields.walletService,
				activeNet:     tt.fields.activeNet,
				walletOpen:    tt.fields.walletOpen,
			}
			if err := c.OpenWallet(); (err != nil) != tt.wantErr {
				t.Errorf("WalletRPCClient.OpenWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWalletRPCClient_CloseWallet(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &WalletRPCClient{
				walletLoader:  tt.fields.walletLoader,
				walletService: tt.fields.walletService,
				activeNet:     tt.fields.activeNet,
				walletOpen:    tt.fields.walletOpen,
			}
			c.CloseWallet()
		})
	}
}

func TestWalletRPCClient_IsWalletOpen(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &WalletRPCClient{
				walletLoader:  tt.fields.walletLoader,
				walletService: tt.fields.walletService,
				activeNet:     tt.fields.activeNet,
				walletOpen:    tt.fields.walletOpen,
			}
			if got := c.IsWalletOpen(); got != tt.want {
				t.Errorf("WalletRPCClient.IsWalletOpen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletRPCClient_SyncBlockChain(t *testing.T) {
	type fields struct {
		walletLoader  walletrpc.WalletLoaderServiceClient
		walletService walletrpc.WalletServiceClient
		activeNet     *netparams.Params
		walletOpen    bool
	}
	tests := []struct {
		name     string
		fields   fields
		listener *app.BlockChainSyncListener
		showLog  bool
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &WalletRPCClient{
				walletLoader:  tt.fields.walletLoader,
				walletService: tt.fields.walletService,
				activeNet:     tt.fields.activeNet,
				walletOpen:    tt.fields.walletOpen,
			}
			if err := c.SyncBlockChain(tt.listener, tt.showLog); (err != nil) != tt.wantErr {
				t.Errorf("WalletRPCClient.SyncBlockChain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

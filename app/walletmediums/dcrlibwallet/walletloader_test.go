package dcrlibwallet

import (
	"testing"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app"
)

func TestDcrWalletLib_NetType(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			if got := lib.NetType(); got != test.want {
				t.Errorf("DcrWalletLib.NetType() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_WalletExists(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.WalletExists()
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.WalletExists() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("DcrWalletLib.WalletExists() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_GenerateNewWalletSeed(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			got, err := lib.GenerateNewWalletSeed()
			if (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.GenerateNewWalletSeed() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("DcrWalletLib.GenerateNewWalletSeed() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_CreateWallet(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			if err := lib.CreateWallet(test.passphrase, test.seed); (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.CreateWallet() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestDcrWalletLib_OpenWallet(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			if err := lib.OpenWallet(); (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.OpenWallet() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestDcrWalletLib_IsWalletOpen(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			if got := lib.IsWalletOpen(); got != test.want {
				t.Errorf("DcrWalletLib.IsWalletOpen() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDcrWalletLib_SyncBlockChain(t *testing.T) {
	type fields struct {
		walletLib *dcrlibwallet.LibWallet
		activeNet *netparams.Params
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lib := &DcrWalletLib{
				walletLib: test.fields.walletLib,
				activeNet: test.fields.activeNet,
			}
			if err := lib.SyncBlockChain(test.listener, test.showLog); (err != nil) != test.wantErr {
				t.Errorf("DcrWalletLib.SyncBlockChain() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

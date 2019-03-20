package config

import (
	"reflect"
	"testing"
)

func TestWalletInfo_MarshalFlag(t *testing.T) {
	type fields struct {
		DbDir   string
		Network string
		Source  string
		Default bool
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
			wallet := &WalletInfo{
				DbDir:   test.fields.DbDir,
				Network: test.fields.Network,
				Source:  test.fields.Source,
				Default: test.fields.Default,
			}
			got, err := wallet.MarshalFlag()
			if (err != nil) != test.wantErr {
				t.Errorf("WalletInfo.MarshalFlag() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("WalletInfo.MarshalFlag() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWalletInfo_UnmarshalFlag(t *testing.T) {
	type fields struct {
		DbDir   string
		Network string
		Source  string
		Default bool
	}
	tests := []struct {
		name    string
		fields  fields
		value   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := &WalletInfo{
				DbDir:   test.fields.DbDir,
				Network: test.fields.Network,
				Source:  test.fields.Source,
				Default: test.fields.Default,
			}
			if err := wallet.UnmarshalFlag(test.value); (err != nil) != test.wantErr {
				t.Errorf("WalletInfo.UnmarshalFlag() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestWalletInfo_Summary(t *testing.T) {
	type fields struct {
		DbDir   string
		Network string
		Source  string
		Default bool
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
			wallet := &WalletInfo{
				DbDir:   test.fields.DbDir,
				Network: test.fields.Network,
				Source:  test.fields.Source,
				Default: test.fields.Default,
			}
			if got := wallet.Summary(); got != test.want {
				t.Errorf("WalletInfo.Summary() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDefaultWallet(t *testing.T) {
	tests := []struct {
		name    string
		wallets []*WalletInfo
		want    *WalletInfo
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := DefaultWallet(test.wallets); !reflect.DeepEqual(got, test.want) {
				t.Errorf("DefaultWallet() = %v, want %v", got, test.want)
			}
		})
	}
}

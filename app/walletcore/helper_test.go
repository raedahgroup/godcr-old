package walletcore

import (
	"testing"

	"github.com/decred/dcrd/chaincfg"
)

func TestGetAddressFromPkScript(t *testing.T) {
	tests := []struct {
		name        string
		activeNet   *chaincfg.Params
		pkScript    []byte
		wantAddress string
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddress, err := GetAddressFromPkScript(tt.activeNet, tt.pkScript)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddressFromPkScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAddress != tt.wantAddress {
				t.Errorf("GetAddressFromPkScript() = %v, want %v", gotAddress, tt.wantAddress)
			}
		})
	}
}

func TestSimpleBalance(t *testing.T) {
	tests := []struct {
		name     string
		balance  *Balance
		detailed bool
		want     string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimpleBalance(tt.balance, tt.detailed); got != tt.want {
				t.Errorf("SimpleBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

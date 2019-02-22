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
		{
			name:      "valid testnet address",
			activeNet: &chaincfg.TestNet3Params,
			pkScript: []byte{0x76, 0xa9, 0x14, 0x29, 0x95, 0xa0,
				0xfe, 0x68, 0x43, 0xfa, 0x9b, 0x95, 0x45,
				0x97, 0xf0, 0xdc, 0xa7, 0xa4, 0x4d, 0xf6,
				0xfa, 0x0b, 0x5c, 0x88, 0xac},
			wantAddress: "TsUp1RpdHGb8wdaNykG49kZwRB3TzVrVDfh",
			wantErr:     false,
		},
		{
			name:      "valid mainnet address",
			activeNet: &chaincfg.MainNetParams,
			pkScript: []byte{0x76, 0xa9, 0x14, 0x29, 0x95, 0xa0,
				0xfe, 0x68, 0x43, 0xfa, 0x9b, 0x95, 0x45,
				0x97, 0xf0, 0xdc, 0xa7, 0xa4, 0x4d, 0xf6,
				0xfa, 0x0b, 0x5c, 0x88, 0xac},
			wantAddress: "DsUknSh7tUY2qGu2AMdu1BYfq55YRkPQKqC",
			wantErr:     false,
		},
		{
			name:      "invalid mainnet address",
			activeNet: &chaincfg.MainNetParams,
			pkScript: []byte{0x76, 0xa9, 0x14, 0x29, 0x95, 0xa0,
				0xfe, 0x68, 0x43, 0xfa, 0x9b, 0x95, 0x45,
				0x97, 0xf0, 0xdc, 0xa7, 0xa4, 0x4d, 0xf6,
				0xfa, 0x0b, 0x5c, 0x88, 0xac},
			wantAddress: "TsUp1RpdHGb8wdaNykG49kZwRB3TzVrVDfh",
			wantErr:     true,
		}
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

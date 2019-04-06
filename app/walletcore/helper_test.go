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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotAddress, err := GetAddressFromPkScript(test.activeNet, test.pkScript)
			if (err != nil) != test.wantErr {
				t.Errorf("GetAddressFromPkScript() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if gotAddress != test.wantAddress {
				t.Errorf("GetAddressFromPkScript() = %v, want %v", gotAddress, test.wantAddress)
			}
		})
	}
}

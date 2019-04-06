package app

import (
	"reflect"
	"testing"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/config"
)

func TestDecredWalletDbDirectories(t *testing.T) {
	tests := []struct {
		name string
		want []WalletDbDir
	}{
		{
			name: "decred wallet db directories",
			want: []WalletDbDir{
				{Source: "dcrwallet", Path: dcrutil.AppDataDir("dcrwallet", false)},
				{Source: "decredition", Path: decreditionAppDirectory()},
				{Source: "godcr", Path: config.DefaultAppDataDir},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := DecredWalletDbDirectories(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("DecredWalletDbDirectories() = %v, want %v", got, test.want)
			}
		})
	}
}

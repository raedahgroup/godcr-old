package app

import (
	"reflect"
	"testing"
)

func TestDecredWalletDbDirectories(t *testing.T) {
	tests := []struct {
		name string
		want []WalletDbDir
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := DecredWalletDbDirectories(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("DecredWalletDbDirectories() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_decreditionAppDirectory(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := decreditionAppDirectory(); got != test.want {
				t.Errorf("decreditionAppDirectory() = %v, want %v", got, test.want)
			}
		})
	}
}

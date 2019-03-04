package dcrlibwallet

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		appDataDir string
		testnet    bool
		want       *DcrWalletLib
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := New(test.appDataDir, test.testnet)
			if (err != nil) != test.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("New() = %v, want %v", got, test.want)
			}
		})
	}
}

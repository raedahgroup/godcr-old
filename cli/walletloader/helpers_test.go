package walletloader

import (
	"context"
	"reflect"
	"testing"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
)

func Test_choseNetworkAndCreateMiddleware(t *testing.T) {
	tests := []struct {
		name    string
		want    app.WalletMiddleware
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := choseNetworkAndCreateMiddleware()
			if (err != nil) != test.wantErr {
				t.Errorf("choseNetworkAndCreateMiddleware() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("choseNetworkAndCreateMiddleware() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_displayWalletSeed(t *testing.T) {
	tests := []struct {
		name string
		seed string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			displayWalletSeed(test.seed)
		})
	}
}

func TestAttemptToCreateWallet(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		want    *config.WalletInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := AttemptToCreateWallet(test.ctx)
			if (err != nil) != test.wantErr {
				t.Errorf("AttemptToCreateWallet() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("AttemptToCreateWallet() = %v, want %v", got, test.want)
			}
		})
	}
}

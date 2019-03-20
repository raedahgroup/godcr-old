package main

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
)

func Test_attemptExecuteSimpleOp(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantIsSimpleOp bool
		wantErr        bool
	}{
		{
			name:           "simple op",
			args:           []string{"cmd", "help", "detect"},
			wantIsSimpleOp: true,
			wantErr:        false,
		},
		{
			name:           "not a simple op",
			args:           []string{"cmd", "balance"},
			wantIsSimpleOp: false,
			wantErr:        false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Args = test.args

			gotIsSimpleOp, err := attemptExecuteSimpleOp()
			if (err != nil) != test.wantErr {
				t.Errorf("attemptExecuteSimpleOp() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if gotIsSimpleOp != test.wantIsSimpleOp {
				t.Errorf("attemptExecuteSimpleOp() = %v, want %v", gotIsSimpleOp, test.wantIsSimpleOp)
			}
		})
	}
}

func Test_connectToWallet(t *testing.T) {
	cfg, _, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	walletMiddleware, err := dcrlibwallet.New(cfg.AppDataDir, config.DefaultWallet(cfg.Wallets))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		ctx     context.Context
		cfg     *config.Config
		want    app.WalletMiddleware
		wantErr bool
	}{
		{
			name:    "valid connection",
			ctx:     context.Background(),
			cfg:     cfg,
			want:    walletMiddleware,
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := connectToWallet(test.ctx, test.cfg)
			if (err != nil) != test.wantErr {
				t.Errorf("connectToWallet() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("connectToWallet() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_detectOrCreateWallet(t *testing.T) {
	conf, _, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	walletInfo := config.DefaultWallet(conf.Wallets)

	tests := []struct {
		name    string
		ctx     context.Context
		want    *config.WalletInfo
		wantErr bool
	}{
		{
			name:    "detect or create wallet",
			ctx:     context.Background(),
			want:    walletInfo,
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := detectOrCreateWallet(test.ctx)
			if (err != nil) != test.wantErr {
				t.Errorf("detectOrCreateWallet() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("detectOrCreateWallet() = %v, want %v", got, test.want)
			}
		})
	}
}

package web

import (
	"context"
	"testing"
	"time"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
	"github.com/raedahgroup/godcr/cli/commands"
)

func TestStartServer(t *testing.T) {
	type test struct {
		name             string
		ctx              context.Context
		cancel           context.CancelFunc
		walletMiddleware app.WalletMiddleware
		host             string
		port             string
		wantErr          bool
	}

	var err error
	var wallets []*config.WalletInfo

	cfg, _, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	wallets = cfg.Wallets
	if wallets == nil {
		wallets, err = commands.DetectWallets(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := make([]test, len(wallets))
	for i := range wallets {
		ctx, cancel := context.WithCancel(ctx)

		walletMiddleware, err := dcrlibwallet.New(cfg.AppDataDir, wallets[i])
		if err != nil {
			t.Fatal(err)
		}
		tests[i] = test{
			name:             "test start server with wallet " + wallets[i].DbDir,
			ctx:              ctx,
			cancel:           cancel,
			walletMiddleware: walletMiddleware,
			host:             cfg.HTTPHost,
			port:             cfg.HTTPPort,
			wantErr:          true,
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			go func() {
				time.AfterFunc(time.Second*2, func() {
					tt.cancel()
				})
			}()

			err := StartServer(tt.ctx, tt.walletMiddleware, tt.host, tt.port)
			if (err != nil) != tt.wantErr && tt.ctx.Err() == nil {
				t.Errorf("StartServer() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func Test_openWalletIfExists(t *testing.T) {
	type test struct {
		name             string
		ctx              context.Context
		cancel           context.CancelFunc
		walletMiddleware app.WalletMiddleware
		wantErr          bool
	}

	var err error
	var wallets []*config.WalletInfo

	cfg, _, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	wallets = cfg.Wallets
	if wallets == nil {
		wallets, err = commands.DetectWallets(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := make([]test, len(wallets))
	for i := range wallets {
		ctx, cancel := context.WithCancel(ctx)

		walletMiddleware, err := dcrlibwallet.New(cfg.AppDataDir, wallets[i])
		if err != nil {
			t.Fatal(err)
		}
		tests[i] = test{
			name:             "open wallet " + wallets[i].DbDir,
			ctx:              ctx,
			cancel:           cancel,
			walletMiddleware: walletMiddleware,
			wantErr:          false,
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				time.AfterFunc(time.Second*2, func() {
					tt.cancel()
				})
			}()

			err := openWalletIfExist(tt.ctx, tt.walletMiddleware)
			if (err != nil) != tt.wantErr && tt.ctx.Err() == nil {
				t.Errorf("openWalletIfExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

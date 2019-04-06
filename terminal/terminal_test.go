package terminal

import (
	"context"
	"testing"
	"time"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
	"github.com/raedahgroup/godcr/cli/commands"
)

func TestStartTerminalApp(t *testing.T) {
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
			name:             "start terminal with wallet " + wallets[i].DbDir,
			ctx:              ctx,
			cancel:           cancel,
			walletMiddleware: walletMiddleware,
			wantErr:          false,
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			go func() {
				time.AfterFunc(time.Second*2, func() {
					test.cancel()
				})
			}()

			err := StartTerminalApp(test.ctx, test.walletMiddleware)
			if (err != nil) != test.wantErr && test.ctx.Err() == nil {
				t.Errorf("StartTerminalApp() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

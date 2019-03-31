package nuklear

import (
	"context"
	"testing"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
	"github.com/raedahgroup/godcr/cli/commands"
)

func TestLaunchApp(t *testing.T) {
	type test struct {
		name             string
		ctx              context.Context
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
		walletMiddleware, err := dcrlibwallet.New(cfg.AppDataDir, wallets[i])
		if err != nil {
			t.Fatal(err)
		}

		tests[i] = test{
			name:             "launch desktop app with wallet " + wallets[i].DbDir,
			ctx:              ctx,
			walletMiddleware: walletMiddleware,
			wantErr:          false,
		}
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := LaunchApp(test.ctx, test.walletMiddleware); (err != nil) != test.wantErr {
				t.Errorf("LaunchApp() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

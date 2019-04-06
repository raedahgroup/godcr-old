package dcrlibwallet

import (
	"reflect"
	"testing"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/util"
	"github.com/raedahgroup/godcr/app/config"
)

func TestNew(t *testing.T) {
	var wallets []*config.WalletInfo
	var err error

	cfg, _, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	wallets = cfg.Wallets

	type test struct {
		name       string
		appDataDir string
		wallet     *config.WalletInfo
		want       *DcrWalletLib
		wantErr    bool
	}

	tests := make([]test, len(wallets))
	for i := range wallets {
		activeNet := util.NetParams(wallets[i].Network)
		if activeNet == nil {
			continue
		}

		lw := dcrlibwallet.LibWalletFromDb(cfg.AppDataDir, wallets[i].DbDir, activeNet)
		lw.SetLogLevel("off")
		lw.InitLoaderWithoutShutdownListener()

		tests[i] = test{
			name:       "new dcrlibwallet " + wallets[i].DbDir,
			appDataDir: cfg.AppDataDir,
			wallet:     wallets[i],
			want: &DcrWalletLib{
				walletLib: lw,
				activeNet: activeNet,
			},
			wantErr: false,
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := New(test.appDataDir, test.wallet)
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

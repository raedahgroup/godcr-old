package commands

import (
	"context"
	"errors"

	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrwalletrpc"
)

func getWalletForTesting() (walletcore.Wallet, error) {
	cfg, _, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	if cfg.WalletRPCServer != "" {
		return dcrwalletrpc.New(context.Background(), cfg.WalletRPCServer, cfg.WalletRPCCert, cfg.NoWalletRPCTLS)
	}

	walletInfo := config.DefaultWallet(cfg.Wallets)
	if walletInfo == nil {
		// no default wallet, ask if to trigger detect command to discover existing wallets or to create new wallet
		walletInfos, err := DetectWallets(context.Background())
		if err != nil {
			return nil, err
		}
		if walletInfos == nil {
			return nil, errors.New("No wallet detected")
		}

		walletInfo = config.DefaultWallet(walletInfos)
	}

	return dcrlibwallet.New(cfg.AppDataDir, walletInfo)
}

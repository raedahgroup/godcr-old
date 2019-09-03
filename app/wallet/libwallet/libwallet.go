package libwallet

import (
	"fmt"
	"context"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/utils"
)

// LibWallet implements `wallet.Wallet` using `dcrlibwallet.LibWallet`
// as medium for connecting to a decred wallet.
type LibWallet struct {
	appCtx      context.Context
	walletDbDir string
	activeNet   *netparams.Params
	dcrlw       *dcrlibwallet.LibWallet
}

func Init(appCtx context.Context, walletDbDir, networkType string) (*LibWallet, error) {
	activeNet := utils.NetParams(networkType)
	if activeNet == nil {
		return nil, fmt.Errorf("unsupported wallet: %s", networkType)
	}

	dcrlibwallet.SetLogLevels("off")
	lw, err := dcrlibwallet.NewLibWalletWithDbPath(walletDbDir, activeNet)
	if err != nil {
		return nil, err
	}

	return &LibWallet{
		appCtx:      appCtx,
		walletDbDir: walletDbDir,
		dcrlw:       lw,
		activeNet:   activeNet,
	}, nil
}

package libwallet

import (
	"fmt"

	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/utils"
)

// LibWallet implements `wallet.Wallet` using `dcrlibwallet.LibWallet`
// as medium for connecting to a decred wallet.
type LibWallet struct {
	walletDbDir string
	activeNet   *netparams.Params
	dcrlw       *dcrlibwallet.LibWallet
}

func Init(walletDbDir, networkType string) (*LibWallet, error) {
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
		walletDbDir: walletDbDir,
		dcrlw:       lw,
		activeNet:   activeNet,
	}, nil
}

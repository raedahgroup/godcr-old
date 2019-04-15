package dcrlibwallet

import (
	"fmt"
	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/utils"
	"github.com/raedahgroup/godcr/app/config"
)

// DcrWalletLib implements `WalletMiddleware` using `dcrlibwallet.LibWallet` as medium for connecting to a decred wallet
// Functions relating to operations that can be performed on a wallet are defined in `walletfunctions.go`
// Other wallet-related functions are defined in `walletloader.go`
type DcrWalletLib struct {
	walletDbDir string
	walletLib   *dcrlibwallet.LibWallet
	activeNet   *netparams.Params
}

// New connects to dcrlibwallet and returns an instance of DcrWalletLib
func New(appDataDir string, wallet *config.WalletInfo) (*DcrWalletLib, error) {
	activeNet := utils.NetParams(wallet.Network)
	if activeNet == nil {
		return nil, fmt.Errorf("unsupported wallet: %s", wallet.Network)
	}

	dcrlibwallet.SetLogLevel("off")
	lw, err := dcrlibwallet.NewLibWalletWithDbPath(wallet.DbDir, activeNet)
	if err != nil {
		return nil, err
	}

	return &DcrWalletLib{
		walletDbDir: wallet.DbDir,
		walletLib:   lw,
		activeNet:   activeNet,
	}, nil
}

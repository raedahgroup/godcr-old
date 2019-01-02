package dcrlibwallet

import (
	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/dcrlibwallet"
)

// MobileWalletLib implements `WalletMiddleware` using `dcrlibwallet.LibWallet` as medium for connecting to a decred wallet
// Functions relating to operations that can be performed on a wallet are defined in `walletfunctions.go`
// Other wallet-related functions are defined in `walletloader.go`
type MobileWalletLib struct {
	walletLib *dcrlibwallet.LibWallet
	activeNet *netparams.Params
}

// New connects to dcrlibwallet and returns an instance of MobileWalletLib
func New(appDataDir string, netType string) *MobileWalletLib {
	lw := dcrlibwallet.NewLibWallet(appDataDir, dcrlibwallet.DefaultDbDriver, netType)
	lw.SetLogLevel("off")
	lw.InitLoaderWithoutShutdownListener()

	var activeNet *netparams.Params
	if netType == "mainnet" {
		activeNet = &netparams.MainNetParams
	} else {
		activeNet = &netparams.TestNet3Params
	}

	return &MobileWalletLib{
		walletLib: lw,
		activeNet: activeNet,
	}
}

package mobilewalletlib

import (
	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/mobilewallet"
)

// MobileWalletLib implements `WalletMiddleware` using `mobilewallet.LibWallet` as medium for connecting to a decred wallet
// Functions relating to operations that can be performed on a wallet are defined in `walletfunctions.go`
// Other wallet-related functions are defined in `walletloader.go`
type MobileWalletLib struct {
	walletLib *mobilewallet.LibWallet
	activeNet *netparams.Params
}

// New connects to mobilewallet and returns an instance of MobileWalletLib
func New(appDataDir string, netType string) *MobileWalletLib {
	lw := mobilewallet.NewLibWallet(appDataDir, mobilewallet.DefaultDbDriver, netType)
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

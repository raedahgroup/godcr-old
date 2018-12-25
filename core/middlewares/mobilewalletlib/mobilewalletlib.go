package mobilewalletlib

import (
	"github.com/decred/dcrwallet/netparams"
	"github.com/raedahgroup/mobilewallet"
)

// MobileWalletClient implements `WalletSource` using `mobilewallet.LibWallet`
// Method implementation of `WalletSource` interface are in functions.go
// Other functions not related to `WalletSource` are in utils.go
type MobileWalletLib struct {
	walletLib *mobilewallet.LibWallet
	activeNet *netparams.Params
}

func New(appDataDir string, netType string) *MobileWalletLib {
	// pass empty db driver to use default
	lw := mobilewallet.NewLibWallet(appDataDir, "", netType)
	lw.SetLogLevel("off")
	lw.InitLoader()

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

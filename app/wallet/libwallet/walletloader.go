package libwallet

import "fmt"

// CreateWallet implements `wallet.Wallet.CreateWallet`.
func (lw *LibWallet) CreateWallet(seed, privatePass string) error {
	return fmt.Errorf("unimplemented")
}

// OpenWallet implements `wallet.Wallet.OpenWallet`.
func (lw *LibWallet) OpenWallet(publicPass string) error {
	return fmt.Errorf("unimplemented")
}

// Shutdown implements `wallet.Wallet.Shutdown`.
func (lw *LibWallet) Shutdown() {

}

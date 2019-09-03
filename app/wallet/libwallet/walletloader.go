package libwallet

import (
	"fmt"
	"time"
)

func (lw *LibWallet) WalletExists() (bool, error) {
	return lw.dcrlw.WalletExists()
}

func (lw *LibWallet) CreateWallet(privatePass, seed string) error {
	return lw.dcrlw.CreateWallet(privatePass, seed)
}

// This method may stall if the wallet database is in use by some other process.
// An error is returned if there's a delay in opening the wallet.
// lw.appCtx is also monitored in the event that a user cancels the operation.
func (lw *LibWallet) OpenWallet(publicPass string) error {
	if lw.dcrlw.WalletOpened() {
		return nil
	}

	// wallet database is opened using bolt db by `github.com/decred/dcrwallet/wallet/internal/bdb`
	// bolt db stalls if the database is currently in use by another process,
	// waiting for the other process to release the file.
	// bold db doc advise setting a 1 second timeout to prevent this stalling.
	// see https://github.com/boltdb/bolt#opening-a-database
	walletOpenDelay := time.NewTicker(5 * time.Second)

	loadWalletDone := make(chan error)
	go func() {
		openWalletError := lw.dcrlw.OpenWallet([]byte(publicPass))
		loadWalletDone <- openWalletError
		walletOpenDelay.Stop()
	}()

	select {
	case <-walletOpenDelay.C:
		return fmt.Errorf("wallet database is in use by another process")

	case err := <-loadWalletDone:
		return err

	case <-lw.appCtx.Done():
		return lw.appCtx.Err()
	}
}

func (lw *LibWallet) ChangePrivatePassphrase(oldPass, newPass string) error {
	return lw.dcrlw.ChangePrivatePassphrase([]byte(oldPass), []byte(newPass))
}

func (lw *LibWallet) NetType() string {
	return lw.activeNet.Params.Name
}

func (lw *LibWallet) Shutdown() {
	lw.dcrlw.Shutdown()
}

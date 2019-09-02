package libwallet

import (
	"context"
	"fmt"
	"time"
)

func (lw *LibWallet) WalletExists() (bool, error) {
	return lw.dcrlw.WalletExists()
}

func (lw *LibWallet) CreateWallet(privatePass, seed string) error {
	return lw.dcrlw.CreateWallet(privatePass, seed)
}

// This method may stall if the wallet database is in use by some other process,
// hence the need for ctx, so user can cancel the operation if it's taking too long
// additionally, an error is returned if there's a delay in opening the wallet.
func (lw *LibWallet) OpenWallet(ctx context.Context, publicPass string) error {
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

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (lw *LibWallet) Shutdown() {
	lw.dcrlw.Shutdown()
}

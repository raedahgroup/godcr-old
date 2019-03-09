package terminal

import (
	"context"
	"fmt"
	"os"

	"github.com/raedahgroup/godcr/app"
	"github.com/rivo/tview"
)

// this method may stall until previous godcr instances are closed (especially in cases of multiple dcrlibwallet instances)
// hence the need for ctx, so user can cancel the operation if it's taking too long
func openWalletIfExist(ctx context.Context, walletMiddleware app.WalletMiddleware) (walletExists bool, err error) {
	var errMsg string
	loadWalletDone := make(chan bool)

	go func() {
		defer func() {
			loadWalletDone <- true
		}()

		walletExists, err = walletMiddleware.WalletExists()
		if err != nil {
			errMsg = fmt.Sprintf("Error checking %s wallet", walletMiddleware.NetType())
		}
		if err != nil || !walletExists {
			return
		}

		err = walletMiddleware.OpenWallet()
		if err != nil {
			errMsg = fmt.Sprintf("Failed to open %s wallet", walletMiddleware.NetType())
		}
	}()

	select {
	case <-loadWalletDone:
		if errMsg != "" {
			fmt.Fprintln(os.Stderr, errMsg)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		return

	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func CreateWallet(tviewApp *tview.Application, seed string, password string, walletMiddleware app.WalletMiddleware) (err error) {
	err = walletMiddleware.CreateWallet(password, seed)
	if err != nil {
		return err
	}

	return nil
}

// syncBlockChain uses the WalletMiddleware provided to download block updates
// this is a long running operation, listen for ctx.Done and stop processing
func SyncBlockChain(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	syncDone := make(chan error)
	go func() {
		syncListener := &app.BlockChainSyncListener{
			SyncStarted: func() {
				fmt.Println("Blockchain sync started")
			},
			SyncEnded: func(err error) {
				if err == nil {
					fmt.Println("Blockchain synced successfully")
				} else {
					fmt.Fprintf(os.Stderr, "Blockchain sync completed with error: %s\n", err.Error())
				}
				syncDone <- err
			},
			OnHeadersFetched:    func(percentageProgress int64) {}, // in cli mode, sync updates are logged to terminal, no need to act on this update alert
			OnDiscoveredAddress: func(state string) {},             // in cli mode, sync updates are logged to terminal, no need to act on update alert
			OnRescanningBlocks:  func(percentageProgress int64) {}, // in cli mode, sync updates are logged to terminal, no need to act on update alert
		}

		err := walletMiddleware.SyncBlockChain(syncListener, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Blockchain sync failed to start. %s\n", err.Error())
			syncDone <- err
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-syncDone:
		return err
	}
}

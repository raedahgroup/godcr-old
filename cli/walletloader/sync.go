package walletloader

import (
	"context"
	"fmt"
	"os"

	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
	"github.com/raedahgroup/godcr/app"
)

// todo review usages
// syncBlockChain uses the WalletMiddleware provided to download block updates
// this is a long running operation, listen for ctx.Done and stop processing
func SyncBlockChain(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	syncError := make(chan error)
	var syncDone bool

	processSyncUpdates := func(report *defaultsynclistener.ProgressReport) {
		if syncDone {
			return
		}

		progressReport := report.Read()

		if progressReport.Done {
			syncDone = true
			if progressReport.Error == "" {
				fmt.Println("Synced successfully.")
				syncError <- nil
			} else {
				fmt.Fprintf(os.Stderr, "Sync completed with error: %s.\n", progressReport.Error)
				syncError <- fmt.Errorf(progressReport.Error)
			}
			return
		}
	}

	fmt.Println("Sync started.")
	walletMiddleware.SpvSync(true, processSyncUpdates)

	// wait for context cancel or sync done trigger before exiting function
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-syncError:
		return err
	}
}

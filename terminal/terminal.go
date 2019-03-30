package terminal

import (
	"context"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/pages"
	"github.com/rivo/tview"
)

func StartTerminalApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	tviewApp := tview.NewApplication()

	walletExists, err := walletMiddleware.OpenWalletIfExist(ctx)
	if err != nil {
		return err
	}

	var page tview.Primitive
	if walletExists {
		page = pages.SyncPage(tviewApp, walletMiddleware)
	} else {
		page = pages.CreateWalletPage(tviewApp, walletMiddleware)
	}

	// `Run` blocks until app.Stop() is called before returning
	return tviewApp.SetRoot(page, true).Run()
}

func sync(tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {
	syncStatus := make(chan walletcore.SyncStatus)
	syncBlockchain(walletMiddleware, syncStatus)

	for status := range syncStatus {
		if status == walletcore.SyncStatusError {
			tviewApp.Stop()
			return nil
		}
		if status == walletcore.SyncStatusInProgress {
			msgOutput := fmt.Sprintf(syncMessage)
			fmt.Println(msgOutput)
			continue
		}
		if status == walletcore.SyncStatusSuccess {
			return pages.TerminalLayout(tviewApp, walletMiddleware)
		}
	}

	return pages.TerminalLayout(tviewApp, walletMiddleware)
}


var syncMessage string

func syncBlockchain(wallet app.WalletMiddleware, syncStatus chan walletcore.SyncStatus) {
	go func() {
		err := wallet.SyncBlockChain(&app.BlockChainSyncListener{
			SyncStarted: func() {
				syncMessage = "Blockchain sync started..."
				syncStatus <- walletcore.SyncStatusInProgress
			},
			SyncEnded: func(err error) {
				if err != nil {
					syncMessage = fmt.Sprintf("Blockchain sync completed with error: %s", err.Error())
					syncStatus <- walletcore.SyncStatusError
				} else {
					syncMessage = "Blockchain sync completed successfully"
					syncStatus <- walletcore.SyncStatusSuccess
				}
			},
			OnHeadersFetched: func(percentageProgress int64) {
				syncMessage = fmt.Sprintf("Blockchain sync in progress. Fetching headers (1/3): %d%%", percentageProgress)
				syncStatus <- walletcore.SyncStatusSuccess
			},
			OnDiscoveredAddress: func(_ string) {
				syncMessage = "Blockchain sync in progress. Discovering addresses (2/3)"
				syncStatus <- walletcore.SyncStatusInProgress
			},
			OnRescanningBlocks: func(percentageProgress int64) {
				syncMessage = fmt.Sprintf("Blockchain sync in progress. Rescanning blocks (3/3): %d%%", percentageProgress)
				syncStatus <- walletcore.SyncStatusInProgress
			},
		}, false)

		if err != nil {
			syncMessage = fmt.Sprintf("Blockchain sync failed to start. %s", err.Error())
			syncStatus <- walletcore.SyncStatusError
		}
	}()
}

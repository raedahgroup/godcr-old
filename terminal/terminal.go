package terminal

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/pages"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func StartTerminalApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	tviewApp := tview.NewApplication()

	walletExists, err := helpers.OpenWalletIfExist(ctx, walletMiddleware)
	if err != nil {
		return err
	}

	var page tview.Primitive
	if walletExists {
		page = pageLoader(tviewApp, walletMiddleware)
	} else {
		page = createWallet(tviewApp, walletMiddleware)
	}

	// `Run` blocks until app.Stop() is called before returning
	return tviewApp.SetRoot(page, true).Run()
}


func pageLoader(tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {
	syncBlockchain(walletMiddleware)
	if Status == walletcore.SyncStatusError {
		msgOutput := fmt.Sprintf(Report)
		helpers.CenterAlignedTextView(msgOutput)
		tviewApp.Stop()
	}
	if Status == walletcore.SyncStatusInProgress {
		msgOutput := fmt.Sprintf(Report)
		return helpers.CenterAlignedTextView(msgOutput)
	}
	if Status == walletcore.SyncStatusSuccess {
		return pages.TerminalLayout(tviewApp, walletMiddleware)
	}

	return nil
}


var Report string
var Status walletcore.SyncStatus

func syncBlockchain(wallet app.WalletMiddleware) {
	err := wallet.SyncBlockChain(&app.BlockChainSyncListener{
		SyncStarted: func() {
			updateStatus("Blockchain sync started...", walletcore.SyncStatusInProgress)
		},
		SyncEnded: func(err error) {
			if err != nil {
				updateStatus(fmt.Sprintf("Blockchain sync completed with error: %s", err.Error()), walletcore.SyncStatusError)
			} else {
				updateStatus("Blockchain sync completed successfully", walletcore.SyncStatusSuccess)
			}
		},
		OnHeadersFetched: func(percentageProgress int64) {
			updateStatus(fmt.Sprintf("Blockchain sync in progress. Fetching headers (1/3): %d%%", percentageProgress), walletcore.SyncStatusInProgress)
		},
		OnDiscoveredAddress: func(_ string) {
			updateStatus("Blockchain sync in progress. Discovering addresses (2/3)", walletcore.SyncStatusInProgress)
		},
		OnRescanningBlocks: func(percentageProgress int64) {
			updateStatus(fmt.Sprintf("Blockchain sync in progress. Rescanning blocks (3/3): %d%%", percentageProgress), walletcore.SyncStatusInProgress)
		},
	}, false)

	if err != nil {
		updateStatus(fmt.Sprintf("Blockchain sync failed to start. %s", err.Error()), walletcore.SyncStatusError)
	}
}

func updateStatus(report string, status walletcore.SyncStatus) {
	Report = report
	Status = status
}

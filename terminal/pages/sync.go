package pages

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func LaunchSyncPage(ctx context.Context, tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) {
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	// page title and hint
	body.AddItem(primitives.NewCenterAlignedTextView("Synchronizing"), 2, 0, false)
	hintText := primitives.WordWrappedTextView("(Press Enter or Esc to cancel sync and exit the app)")
	hintText.SetTextColor(tcell.ColorGray)
	body.AddItem(hintText, 3, 0, false)

	// text view to show sync progress updates
	syncStatusTextView := primitives.NewCenterAlignedTextView("")
	body.AddItem(syncStatusTextView, 3, 0, false)

	cancelButton := tview.NewButton("Cancel and Exit")
	body.AddItem(cancelButton, 1, 0, true)

	var appStopped bool
	cancelAndExit := func() {
		if appStopped {
			return
		}
		appStopped = true
		tviewApp.Stop()
	}
	cancelButton.SetSelectedFunc(cancelAndExit)
	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			cancelAndExit()
			return nil
		}
		return event
	})

	tviewApp.SetRoot(body, true)

	// function to update the sync page with status report from the sync operation
	updateStatus := func(status string) {
		if appStopped {
			return
		}
		tviewApp.QueueUpdateDraw(func() {
			syncStatusTextView.SetText(status)
		})
	}

	// function to be executed after the sync operation completes successfully
	afterSyncing := func() {
		if appStopped {
			return
		}
		tviewApp.QueueUpdateDraw(func() {
			tviewApp.SetRoot(rootPage(tviewApp, walletMiddleware), true)
		})
	}

	startSync(ctx, walletMiddleware, updateStatus, afterSyncing)
}

// this is a long running operation, listen for ctx.Done and stop processing
func startSync(_ context.Context, walletMiddleware app.WalletMiddleware, updateStatus func(string), afterSyncing func()) {
	err := walletMiddleware.SyncBlockChain(&app.BlockChainSyncListener{
		SyncStarted: func() {
			updateStatus("Blockchain sync started...")
		},
		SyncEnded: func(err error) {
			if err != nil {
				updateStatus(fmt.Sprintf("Blockchain sync completed with error: %s", err.Error()))
			} else {
				updateStatus("Blockchain sync completed successfully")
				afterSyncing()
			}
		},
		OnHeadersFetched: func(percentageProgress int64) {
			updateStatus(fmt.Sprintf("Blockchain sync in progress. Fetching headers (1/3): %d%%", percentageProgress))
		},
		OnDiscoveredAddress: func(_ string) {
			updateStatus("Blockchain sync in progress. Discovering addresses (2/3)")
		},
		OnRescanningBlocks: func(percentageProgress int64) {
			updateStatus(fmt.Sprintf("Blockchain sync in progress. Rescanning blocks (3/3): %d%%", percentageProgress))
		},
	}, false)

	if err != nil {
		updateStatus(fmt.Sprintf("Blockchain sync failed to start. %s", err.Error()))
	}
}

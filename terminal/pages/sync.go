package pages

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
	"github.com/raedahgroup/godcr/app/sync"
	"time"
)

func LaunchSyncPage(tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) {
	syncPage := tview.NewFlex().SetDirection(tview.FlexRow)

	// page title and hint
	syncPage.AddItem(primitives.NewCenterAlignedTextView("Synchronizing"), 1, 0, false)
	hintText := primitives.WordWrappedTextView("(Press ESC twice to cancel sync and exit the app)").SetTextColor(helpers.HintTextColor)
	syncPage.AddItem(hintText, 3, 0, false)

	errorTextView := primitives.WordWrappedTextView("")
	errorTextView.SetTextColor(helpers.DecredOrangeColor)

	// function to display sync errors
	handleError := func(errorMessage string) {
		tviewApp.QueueUpdateDraw(func() {
			syncPage.RemoveItem(errorTextView)
			errorTextView.SetText(errorMessage)
			syncPage.AddItem(errorTextView, 3, 0, false)
		})
	}

	// function to update the sync page with status report from the sync operation
	var previousUpdateViews []*primitives.TextView
	updateStatus := func(report []string) {
		tviewApp.QueueUpdateDraw(func() {
			// remove previous update views and error view
			for _, view := range previousUpdateViews {
				syncPage.RemoveItem(view)
			}
			syncPage.RemoveItem(errorTextView)
			previousUpdateViews = make([]*primitives.TextView, len(report))

			for i, info := range report {
				previousUpdateViews[i] = primitives.NewCenterAlignedTextView(info)
				syncPage.AddItem(previousUpdateViews[i], 1, 0, false)
			}

			// re-display error view?
			if errorTextView.GetText() != "" {
				syncPage.AddItem(errorTextView, 3, 0, false)
			}
		})
	}

	// function to be executed after the sync operation completes successfully
	afterSyncing := func() {
		tviewApp.QueueUpdateDraw(func() {
			tviewApp.SetRoot(rootPage(tviewApp, walletMiddleware), true)
		})
	}

	startSync(walletMiddleware, updateStatus, handleError, afterSyncing)

	var cancelTriggered bool
	syncPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			if cancelTriggered {
				tviewApp.Stop()
			} else {
				cancelTriggered = true
				// remove cancel trigger after 1 second if user does not press escape again within that time
				go func() {
					<- time.After(1 * time.Second)
					cancelTriggered = false
				}()
			}
			return nil
		}
		return event
	})

	syncPage.SetBorderPadding(3, 3, 3, 3)
	syncPage.SetBackgroundColor(tcell.ColorBlack)
	tviewApp.SetRoot(syncPage, true)
}

func startSync(walletMiddleware app.WalletMiddleware, updateStatus func([]string), handleError func(string), afterSyncing func()) {
	err := walletMiddleware.SyncBlockChain(false, func(syncPrivateInfo *sync.PrivateInfo) {
		syncInfo := syncPrivateInfo.Read()
		if syncInfo.Status == sync.StatusSuccess {
			afterSyncing()
			return
		}

		var report []string
		if syncInfo.TotalTimeRemaining == "" {
			report = []string{
				fmt.Sprintf("%d%% completed.", syncInfo.TotalSyncProgress),
			}
		} else {
			report = []string{
				fmt.Sprintf("%d%% completed, %s remaining.", syncInfo.TotalSyncProgress, syncInfo.TotalTimeRemaining),
			}
		}

		switch syncInfo.CurrentStep {
		case 1:
			report = append(report, fmt.Sprintf("Fetched %d of %d block headers.",
				syncInfo.FetchedHeadersCount, syncInfo.TotalHeadersToFetch))
			report = append(report, fmt.Sprintf("%d%% through step 1 of 3.", syncInfo.HeadersFetchProgress))

			if syncInfo.DaysBehind != "" {
				report = append(report, fmt.Sprintf("Your wallet is %s behind.", syncInfo.DaysBehind))
			}

		case 2:
			report = append(report, "Discovering used addresses.")
			if syncInfo.AddressDiscoveryProgress > 100 {
				report = append(report, fmt.Sprintf("%d%% (over) through step 2 of 3.", syncInfo.AddressDiscoveryProgress))
			} else {
				report = append(report, fmt.Sprintf("%d%% through step 2 of 3.", syncInfo.AddressDiscoveryProgress))
			}

		case 3:
			report = append(report, fmt.Sprintf("Scanning %d of %d block headers.",
				syncInfo.CurrentRescanHeight, syncInfo.TotalHeadersToFetch))
			report = append(report, fmt.Sprintf("%d%% through step 3 of 3.", syncInfo.HeadersFetchProgress))
		}

		// show peer count last
		if syncInfo.ConnectedPeers == 1 {
			report = append(report, fmt.Sprintf("Syncing with %d peer on %s", syncInfo.ConnectedPeers, walletMiddleware.NetType()))
		} else {
			report = append(report, fmt.Sprintf("Syncing with %d peers on %s", syncInfo.ConnectedPeers, walletMiddleware.NetType()))
		}

		updateStatus(report)

		if syncInfo.Status == sync.StatusError {
			handleError("Sync error: "+syncInfo.Error)
		}
	})

	if err != nil {
		handleError(fmt.Sprintf("Sync failed to start: %s", err.Error()))
	} else {
		updateStatus([]string{"Starting..."})
	}
}

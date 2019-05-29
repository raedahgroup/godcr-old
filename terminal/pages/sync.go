package pages

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func LaunchSyncPage(tviewApp *tview.Application, walletMiddleware app.WalletMiddleware, displayPage func(tview.Primitive), hintTextView *primitives.TextView, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	syncPage := tview.NewFlex().SetDirection(tview.FlexRow)

	// page title
	syncPage.AddItem(primitives.NewCenterAlignedTextView("Synchronizing"), 1, 0, false)

	errorTextView := primitives.WordWrappedTextView("")
	errorTextView.SetTextColor(helpers.DecredOrangeColor)

	// function to display sync errors
	handleError := func(errorMessage string) {
		tviewApp.QueueUpdateDraw(func() {
			syncPage.RemoveItem(errorTextView)
			errorTextView.SetText(errorMessage)
			syncPage.AddItem(errorTextView, 1, 0, false)
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
				syncPage.AddItem(errorTextView, 1, 0, false)
			}
		})
	}

	// function to be executed after the sync operation completes successfully
	afterSyncing := func() {
		tviewApp.QueueUpdateDraw(func() {
			clearFocus()
			displayPage(overviewPage(walletMiddleware, hintTextView, tviewApp, clearFocus))
		})
	}

	startSync(walletMiddleware, updateStatus, handleError, afterSyncing)

	var cancelTriggered bool
	syncPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			if cancelTriggered {
				clearFocus()
				displayPage(overviewPage(walletMiddleware, hintTextView, tviewApp, clearFocus))
			} else {
				cancelTriggered = true
				// remove cancel trigger after 1 second if user does not press escape again within that time
				go func() {
					<-time.After(1 * time.Second)
					cancelTriggered = false
				}()
			}
			return nil
		}
		return event
	})

	hintTextView.SetText("Press ESC twice to cancel synchronization process").SetTextColor(helpers.HintTextColor)

	setFocus(syncPage)
	return syncPage
}

func startSync(walletMiddleware app.WalletMiddleware, updateStatus func([]string), handleError func(string), afterSyncing func()) {
	walletMiddleware.SyncBlockChain(false, func(report *defaultsynclistener.ProgressReport) {
		progressReport := report.Read()

		if progressReport.Status == defaultsynclistener.SyncStatusSuccess {
			afterSyncing()
			return
		}

		var stringReport []string
		if progressReport.TotalTimeRemaining == "" {
			stringReport = []string{
				fmt.Sprintf("%d%% completed.", progressReport.TotalSyncProgress),
			}
		} else {
			stringReport = []string{
				fmt.Sprintf("%d%% completed, %s remaining.", progressReport.TotalSyncProgress, progressReport.TotalTimeRemaining),
			}
		}

		switch progressReport.CurrentStep {
		case defaultsynclistener.FetchingBlockHeaders:
			stringReport = append(stringReport, fmt.Sprintf("Fetched %d of %d block headers.",
				progressReport.FetchedHeadersCount, progressReport.TotalHeadersToFetch))
			stringReport = append(stringReport, fmt.Sprintf("%d%% through step 1 of 3.", progressReport.HeadersFetchProgress))

			if progressReport.DaysBehind != "" {
				stringReport = append(stringReport, fmt.Sprintf("Your wallet is %s behind.", progressReport.DaysBehind))
			}

		case defaultsynclistener.DiscoveringUsedAddresses:
			stringReport = append(stringReport, "Discovering used addresses.")
			if progressReport.AddressDiscoveryProgress > 100 {
				stringReport = append(stringReport, fmt.Sprintf("%d%% (over) through step 2 of 3.", progressReport.AddressDiscoveryProgress))
			} else {
				stringReport = append(stringReport, fmt.Sprintf("%d%% through step 2 of 3.", progressReport.AddressDiscoveryProgress))
			}

		case defaultsynclistener.ScanningBlockHeaders:
			stringReport = append(stringReport, fmt.Sprintf("Scanning %d of %d block headers.",
				progressReport.CurrentRescanHeight, progressReport.TotalHeadersToFetch))
			stringReport = append(stringReport, fmt.Sprintf("%d%% through step 3 of 3.", progressReport.HeadersFetchProgress))
		}

		// show peer count last
		if progressReport.ConnectedPeers == 1 {
			stringReport = append(stringReport, fmt.Sprintf("Syncing with %d peer on %s.", progressReport.ConnectedPeers, walletMiddleware.NetType()))
		} else {
			stringReport = append(stringReport, fmt.Sprintf("Syncing with %d peers on %s.", progressReport.ConnectedPeers, walletMiddleware.NetType()))
		}

		updateStatus(stringReport)

		if progressReport.Status == defaultsynclistener.SyncStatusError {
			handleError("Sync error: " + progressReport.Error)
		}
	})
}

package pages

import (
	"fmt"

	"github.com/decred/dcrd/dcrutil"
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func overviewPage() tview.Primitive {
	overviewPage := tview.NewFlex().SetDirection(tview.FlexRow)

	renderBalance(overviewPage)

	// single line space between balance and recent activity section
	overviewPage.AddItem(nil, 1, 0, false)

	renderRecentActivity(overviewPage)

	// single line space between recent activity and sync section
	overviewPage.AddItem(nil, 1, 0, false)

	renderSyncStatus(overviewPage)

	commonPageData.hintTextView.SetText("TIP: Scroll recent activity table with ARROW KEYS. Return to navigation menu with ESC")

	overviewPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			commonPageData.clearAllPageContent()
			return nil
		}
		return event
	})

	commonPageData.app.SetFocus(overviewPage)

	return overviewPage
}

func renderBalance(overviewPage *tview.Flex) {
	balanceTitleTextView := primitives.NewLeftAlignedTextView("Balance")
	overviewPage.AddItem(balanceTitleTextView, 2, 0, false)

	accounts, err := commonPageData.wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		errorTextView := primitives.NewCenterAlignedTextView(err.Error()).SetTextColor(helpers.DecredOrangeColor)
		overviewPage.AddItem(errorTextView, 2, 0, false)
		return
	}

	var totalBalance, spendableBalance dcrutil.Amount
	for _, account := range accounts.Acc {
		totalBalance += dcrutil.Amount(account.Balance.Total)
		spendableBalance += dcrutil.Amount(account.Balance.Total)
	}

	var balance string
	if totalBalance != spendableBalance {
		balance = fmt.Sprintf("Total %s (Spendable %s)", totalBalance.String(), spendableBalance.String())
	} else {
		balance = totalBalance.String()
	}

	balanceTextView := primitives.NewLeftAlignedTextView(balance)
	overviewPage.AddItem(balanceTextView, 2, 0, false)
}

func renderRecentActivity(overviewPage *tview.Flex) {
	overviewPage.AddItem(primitives.NewLeftAlignedTextView("-Recent Activity-").SetTextColor(helpers.DecredLightBlueColor), 1, 0, false)

	statusTextView := primitives.NewCenterAlignedTextView("")
	displayMessage := func(message string, error bool) {
		overviewPage.RemoveItem(statusTextView)
		if message != "" {
			if error {
				statusTextView.SetTextColor(helpers.DecredOrangeColor)
			} else {
				statusTextView.SetTextColor(tcell.ColorWhite)
			}

			statusTextView.SetText(message)
			overviewPage.AddItem(statusTextView, 2, 0, false)
		}
	}

	displayMessage("Fetching data...", false)

	txns, err := commonPageData.wallet.GetTransactionsRaw(0, 5, dcrlibwallet.TxFilterAll)
	if err != nil {
		// updating an element on the page from a goroutine, use tviewApp.QueueUpdateDraw
		commonPageData.app.QueueUpdateDraw(func() {
			displayMessage(err.Error(), true)
		})
		return
	}

	if len(txns) == 0 {
		displayMessage("No activity yet", false)
		return
	}

	historyTable := primitives.NewTable()
	historyTable.SetBorders(false).SetFixed(1, 0)

	// historyTable header
	historyTable.SetHeaderCell(0, 0, "Date (UTC)")
	historyTable.SetHeaderCell(0, 1, fmt.Sprintf("%10s", "Direction"))
	historyTable.SetHeaderCell(0, 2, fmt.Sprintf("%8s", "Amount"))
	historyTable.SetHeaderCell(0, 3, fmt.Sprintf("%5s", "Status"))
	historyTable.SetHeaderCell(0, 4, fmt.Sprintf("%-5s", "Type"))

	// calculate max number of digits after decimal point for all amounts for 5 most recent txs
	inputsAndOutputsAmount := make([]int64, 5)
	for i, tx := range txns {
		if i < 5 {
			inputsAndOutputsAmount[i] = tx.Amount
		} else {
			break
		}
	}
	maxDecimalPlacesForTxAmounts := helpers.MaxDecimalPlaces(inputsAndOutputsAmount)

	for _, tx := range txns {
		nextRowIndex := historyTable.GetRowCount()

		dateCell := tview.NewTableCell(fmt.Sprintf("%-10s", dcrlibwallet.ExtractDateOrTime(tx.Timestamp))).
			SetAlign(tview.AlignCenter).
			SetMaxWidth(1).
			SetExpansion(1)
		historyTable.SetCell(nextRowIndex, 0, dateCell)

		directionCell := tview.NewTableCell(fmt.Sprintf("%-10s", dcrlibwallet.TransactionDirectionName(tx.Direction))).
			SetAlign(tview.AlignCenter).
			SetMaxWidth(2).
			SetExpansion(1)
		historyTable.SetCell(nextRowIndex, 1, directionCell)

		formattedAmount := helpers.FormatAmountDisplay(tx.Amount, maxDecimalPlacesForTxAmounts)
		amountCell := tview.NewTableCell(fmt.Sprintf("%15s", formattedAmount)).
			SetAlign(tview.AlignCenter).
			SetMaxWidth(3).
			SetExpansion(1)
		historyTable.SetCell(nextRowIndex, 2, amountCell)

		status := "Pending"
		confirmations := commonPageData.wallet.GetBestBlock() - tx.BlockHeight + 1
		if tx.BlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
			status = "Confirmed"
		}
		statusCell := tview.NewTableCell(fmt.Sprintf("%12s", status)).
			SetAlign(tview.AlignCenter).
			SetMaxWidth(1).
			SetExpansion(1)
		historyTable.SetCell(nextRowIndex, 3, statusCell)

		typeCell := tview.NewTableCell(fmt.Sprintf("%-8s", tx.Type)).
			SetAlign(tview.AlignCenter).
			SetMaxWidth(1).
			SetExpansion(1)
		historyTable.SetCell(nextRowIndex, 4, typeCell)
	}

	historyTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			commonPageData.clearAllPageContent()
		}
	})

	overviewPage.RemoveItem(statusTextView)
	overviewPage.AddItem(historyTable, 0, 1, true)
	commonPageData.app.SetFocus(historyTable)
}

func renderSyncStatus(overviewPage *tview.Flex) {
	overviewPage.AddItem(primitives.NewLeftAlignedTextView("-Sync Status-").SetTextColor(helpers.DecredLightBlueColor), 1, 0, false)

	if commonPageData.wallet.IsSynced() {
		overviewPage.AddItem(primitives.NewCenterAlignedTextView("Wallet is synced."), 1, 0, false)
		return
	}

	(&syncProgressListener{overviewPage: overviewPage}).start()
}

type syncProgressListener struct {
	overviewPage      *tview.Flex
	updateViews       []*primitives.TextView
	peerCountTextView *primitives.TextView
	errorTextView     *primitives.TextView
}

func (listener *syncProgressListener) start() {
	listener.peerCountTextView = primitives.NewCenterAlignedTextView("")
	listener.errorTextView = primitives.WordWrappedTextView("")

	listener.overviewPage.AddItem(listener.peerCountTextView, 1, 0, false)

	commonPageData.wallet.AddSyncProgressListener(listener, "terminal-ui")
}

func (listener *syncProgressListener) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	commonPageData.app.QueueUpdateDraw(func() {
		peerConnectionSummary := fmt.Sprintf("Syncing with %d peer on %s.", numberOfConnectedPeers,
			commonPageData.wallet.NetType())
		listener.peerCountTextView.SetText(peerConnectionSummary)
	})
}

func (listener *syncProgressListener) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {
	var report = []string{
		fmt.Sprintf("%d%% completed, %s remaining.", headersFetchProgress.TotalSyncProgress,
			dcrlibwallet.CalculateTotalTimeRemaining(headersFetchProgress.TotalTimeRemainingSeconds)),

		fmt.Sprintf("Fetched %d of %d block headers.", headersFetchProgress.FetchedHeadersCount,
			headersFetchProgress.TotalHeadersToFetch),

		fmt.Sprintf("%d%% through step 1 of 3.", headersFetchProgress.HeadersFetchProgress),

		fmt.Sprintf("Your wallet is %s behind.",
			dcrlibwallet.CalculateDaysBehind(headersFetchProgress.CurrentHeaderTimestamp)),
	}
	listener.updateUI(report)
}

func (listener *syncProgressListener) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {
	var report = []string{
		fmt.Sprintf("%d%% completed, %s remaining.", addressDiscoveryProgress.TotalSyncProgress,
			dcrlibwallet.CalculateTotalTimeRemaining(addressDiscoveryProgress.TotalTimeRemainingSeconds)),

		"%Discovering used addresses.",
	}

	if addressDiscoveryProgress.AddressDiscoveryProgress > 100 {
		report = append(report, fmt.Sprintf("%d%% (over) through step 2 of 3.",
			addressDiscoveryProgress.AddressDiscoveryProgress))
	} else {
		report = append(report, fmt.Sprintf("%d%% through step 2 of 3.",
			addressDiscoveryProgress.AddressDiscoveryProgress))
	}

	listener.updateUI(report)
}

func (listener *syncProgressListener) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {
	var report = []string{
		fmt.Sprintf("%d%% completed, %s remaining.", headersRescanProgress.TotalSyncProgress,
			dcrlibwallet.CalculateTotalTimeRemaining(headersRescanProgress.TotalTimeRemainingSeconds)),

		fmt.Sprintf("Scanning %d of %d block headers.", headersRescanProgress.CurrentRescanHeight,
			headersRescanProgress.TotalHeadersToScan),

		fmt.Sprintf("%d%% through step 3 of 3.", headersRescanProgress.RescanProgress),
	}
	listener.updateUI(report)
}

func (listener *syncProgressListener) OnSyncCompleted() {
	// remove previous update views and error view
	for _, view := range listener.updateViews {
		listener.overviewPage.RemoveItem(view)
	}
	listener.overviewPage.RemoveItem(listener.peerCountTextView)
	listener.overviewPage.RemoveItem(listener.errorTextView)

	listener.overviewPage.AddItem(primitives.NewCenterAlignedTextView("Wallet is synced."), 1, 0, false)
}

func (listener *syncProgressListener) OnSyncCanceled(willRestart bool) {}

func (listener *syncProgressListener) OnSyncEndedWithError(err error) {
	// remove previous update views and error view
	for _, view := range listener.updateViews {
		listener.overviewPage.RemoveItem(view)
	}
	listener.overviewPage.RemoveItem(listener.peerCountTextView)
	listener.overviewPage.RemoveItem(listener.errorTextView)

	listener.overviewPage.AddItem(listener.errorTextView, 1, 0, false)
	listener.errorTextView.SetText(fmt.Sprintf("Sync error: %v", err))
}

func (listener *syncProgressListener) Debug(debugInfo *dcrlibwallet.DebugInfo) {}

func (listener *syncProgressListener) updateUI(report []string) {
	commonPageData.app.QueueUpdateDraw(func() {
		// remove previous update views and error view
		for _, view := range listener.updateViews {
			listener.overviewPage.RemoveItem(view)
		}
		listener.overviewPage.RemoveItem(listener.peerCountTextView)
		listener.overviewPage.RemoveItem(listener.errorTextView)

		listener.updateViews = make([]*primitives.TextView, len(report))
		for i, info := range report {
			listener.updateViews[i] = primitives.NewCenterAlignedTextView(info)
			listener.overviewPage.AddItem(listener.updateViews[i], 1, 0, false)
		}

		listener.overviewPage.AddItem(listener.peerCountTextView, 1, 0, false)

		// re-display error view?
		if listener.errorTextView.GetText() != "" {
			listener.overviewPage.AddItem(listener.errorTextView, 1, 0, false)
		}
	})
}

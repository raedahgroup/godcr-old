package pages

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/dcrlibwallet/utils"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func overviewPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, tviewApp *tview.Application, clearFocus func()) tview.Primitive {
	overviewPage := tview.NewFlex().SetDirection(tview.FlexRow)

	renderBalance(overviewPage, wallet)

	// single line space between balance and recent activity section
	overviewPage.AddItem(nil, 1, 0, false)

	// fetch recent activity in subroutine, so the UI doesn't become unresponsive
	go renderRecentActivity(overviewPage, wallet, tviewApp, clearFocus)

	hintTextView.SetText("TIP: Scroll recent activity table with ARROW KEYS. Return to navigation menu with ESC")

	overviewPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			clearFocus()
			return nil
		}

		return event
	})

	tviewApp.SetFocus(overviewPage)

	return overviewPage
}

func renderBalance(overviewPage *tview.Flex, wallet walletcore.Wallet) {
	balanceTitleTextView := primitives.NewLeftAlignedTextView("Balance")
	overviewPage.AddItem(balanceTitleTextView, 2, 0, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		errorTextView := primitives.NewCenterAlignedTextView(err.Error()).SetTextColor(helpers.DecredOrangeColor)
		overviewPage.AddItem(errorTextView, 2, 0, false)
	} else {
		balanceTextView := primitives.NewLeftAlignedTextView(walletcore.WalletBalance(accounts))
		overviewPage.AddItem(balanceTextView, 2, 0, false)
	}
}

func renderRecentActivity(overviewPage *tview.Flex, wallet walletcore.Wallet, tviewApp *tview.Application, clearFocus func()) {
	overviewPage.AddItem(primitives.NewLeftAlignedTextView("-Recent Activity-").SetTextColor(helpers.DecredLightBlueColor), 1, 0, false)

	statusTextView := primitives.NewCenterAlignedTextView("").SetTextColor(helpers.DecredOrangeColor)
	// adding an element to the page from a goroutine, use tviewApp.QueueUpdateDraw
	tviewApp.QueueUpdateDraw(func() {
		overviewPage.AddItem(statusTextView, 2, 0, false)
	})

	txns, err := wallet.TransactionHistory(0, 5, nil)
	if err != nil {
		// updating an element on the page from a goroutine, use tviewApp.QueueUpdateDraw
		tviewApp.QueueUpdateDraw(func() {
			statusTextView.SetText(err.Error())
		})
		return
	}
	if len(txns) == 0 {
		noTxnsTextview := primitives.NewCenterAlignedTextView("No activity yet")
		tviewApp.QueueUpdateDraw(func() {
			overviewPage.AddItem(noTxnsTextview, 2, 0, false)
		})
		return
	}

	statusTextView.SetText("Fetching data...")

	historyTable := primitives.NewTable()
	historyTable.SetBorders(false).SetFixed(1, 0)

	// historyTable header
	historyTable.SetHeaderCell(0, 0, "Date (UTC)")
	historyTable.SetHeaderCell(0, 1, (fmt.Sprintf("%10s", "Direction")))
	historyTable.SetHeaderCell(0, 2, (fmt.Sprintf("%8s", "Amount")))
	historyTable.SetHeaderCell(0, 3, (fmt.Sprintf("%5s", "Status")))
	historyTable.SetHeaderCell(0, 4, (fmt.Sprintf("%-5s", "Type")))

	// calculate max number of digits after decimal point for all amounts for 5 most recent txs
	inputsAndOutputsAmount := make([]int64, 5)
	for i, tx := range txns {
		if i < 5 {
			inputsAndOutputsAmount[i] = tx.Amount
		} else {
			break
		}
	}
	maxDecimalPlacesForTxAmounts := maxDecimalPlaces(inputsAndOutputsAmount)

	// now format amount having determined the max number of decimal places
	formatAmount := func(amount int64) string {
		return formatAmountDisplay(amount, maxDecimalPlacesForTxAmounts)
	}

	for _, tx := range txns {
		nextRowIndex := historyTable.GetRowCount()
		if nextRowIndex > 5 {
			break
		}

		historyTable.SetCell(nextRowIndex, 0, tview.NewTableCell(fmt.Sprintf("%-10s", utils.ExtractDateOrTime(tx.Timestamp))).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1).SetMaxWidth(1).SetExpansion(1))
		historyTable.SetCell(nextRowIndex, 1, tview.NewTableCell(fmt.Sprintf("%-10s", tx.Direction.String())).SetAlign(tview.AlignCenter).SetMaxWidth(2).SetExpansion(1))
		historyTable.SetCell(nextRowIndex, 2, tview.NewTableCell(fmt.Sprintf("%15s", formatAmount(tx.Amount))).SetAlign(tview.AlignCenter).SetMaxWidth(3).SetExpansion(1))
		historyTable.SetCell(nextRowIndex, 3, tview.NewTableCell(fmt.Sprintf("%12s", tx.Status)).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1))
		historyTable.SetCell(nextRowIndex, 4, tview.NewTableCell(fmt.Sprintf("%-8s", tx.Type)).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1))
	}

	historyTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			clearFocus()
		}
	})

	// adding an element to the page from a goroutine, use tviewApp.QueueUpdateDraw
	tviewApp.QueueUpdateDraw(func() {
		overviewPage.RemoveItem(statusTextView)
		overviewPage.AddItem(historyTable, 0, 1, true)
		tviewApp.SetFocus(historyTable)
	})

	return
}

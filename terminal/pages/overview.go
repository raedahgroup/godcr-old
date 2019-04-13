package pages

import (
	"context"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func overviewPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	overviewPage := tview.NewFlex().SetDirection(tview.FlexRow)

	renderBalance(overviewPage, wallet)

	// single line space between balance and recent activity section
	overviewPage.AddItem(nil, 1, 0, false)

	renderRecentActivity(overviewPage, wallet, clearFocus)

	hintTextView.SetText("TIP: Scroll recent activity table with ARROW KEYS. Return to navigation menu with ESC")

	setFocus(overviewPage)

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

func renderRecentActivity(overviewPage *tview.Flex, wallet walletcore.Wallet, clearFocus func()) {
	overviewPage.AddItem(primitives.NewLeftAlignedTextView("-Recent Activity-").SetTextColor(helpers.DecredLightBlueColor), 1, 0, false)

	txns, _, err := wallet.TransactionHistory(context.Background(), -1, 5)
	if err != nil {
		overviewPage.AddItem(primitives.NewCenterAlignedTextView(err.Error()), 2, 0, false)
		return
	}

	historyTable := primitives.NewTable()
	historyTable.SetBorders(false).SetFixed(1, 0)

	// historyTable header
	historyTable.SetHeaderCell(0, 0, "Date (UTC)")
	historyTable.SetHeaderCell(0, 1, "Direction")
	historyTable.SetHeaderCell(0, 2, "Amount")
	historyTable.SetHeaderCell(0, 3, "Status")
	historyTable.SetHeaderCell(0, 4, "Type")

	for _, tx := range txns {
		row := historyTable.GetRowCount()
		if row >= 5 {
			break
		}

		historyTable.SetCell(row, 0, tview.NewTableCell(tx.FormattedTime).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 1, tview.NewTableCell(tx.Direction.String()).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 2, tview.NewTableCell(tx.Amount).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 3, tview.NewTableCell(tx.Status).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 4, tview.NewTableCell(tx.Type).SetAlign(tview.AlignCenter))
	}

	historyTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			clearFocus()
		}
	})

	overviewPage.AddItem(historyTable, 0, 1, true)
	return
}

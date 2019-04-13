package pages

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

var _24hours, _ = time.ParseDuration("24h")

func overviewPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	overviewPage := tview.NewFlex().SetDirection(tview.FlexRow)

	errorTextView := primitives.WordWrappedTextView("")
	errorTextView.SetTextColor(helpers.DecredOrangeColor)

	displayError := func(errorMessage string) {
		overviewPage.RemoveItem(errorTextView)
		errorTextView.SetText(errorMessage)
		overviewPage.AddItem(errorTextView, 2, 0, false)
	}

	balanceTitleTextView := primitives.NewLeftAlignedTextView("Balance")
	overviewPage.AddItem(balanceTitleTextView, 2, 0, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return overviewPage.AddItem(primitives.NewCenterAlignedTextView(err.Error()).SetTextColor(helpers.DecredOrangeColor), 2, 0, false)
	}

	overviewPage.AddItem(primitives.NewLeftAlignedTextView(walletcore.WalletBalance(accounts)), 2, 0, false)

	overviewPage.AddItem(nil, 1, 0, false) // em

	renderRecentActivity(overviewPage, wallet, displayError, clearFocus)

	hintTextView.SetText("TIP: Scroll recent activity table with ARROW KEYS. Return to navigation menu with ESC")

	setFocus(overviewPage)

	return overviewPage
}

func renderRecentActivity(overviewPage *tview.Flex, wallet walletcore.Wallet, displayError func(message string), clearFocus func()) {
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

	loc, _ := time.LoadLocation("UTC")
	currentDate := time.Now().In(loc)

	for _, tx := range txns {
		row := historyTable.GetRowCount()
		if row >= 5 {
			break
		}
		txns, err := wallet.GetTransaction(tx.Hash)
		if err != nil {
			displayError(err.Error())
			return
		}

		transactionDate := time.Unix(tx.Timestamp, 0).In(loc)
		transactionAge := currentDate.Sub(transactionDate)
		txDateTime := strings.Split(tx.FormattedTime, " ")

		if transactionAge > _24hours {
			historyTable.SetCell(row, 0, tview.NewTableCell(fmt.Sprintln(txDateTime[0])).SetAlign(tview.AlignCenter))
		} else {
			historyTable.SetCell(row, 0, tview.NewTableCell(fmt.Sprintln(txDateTime[1])).SetAlign(tview.AlignCenter))
		}

		if txns.Confirmations >= walletcore.DefaultRequiredConfirmations {
			historyTable.SetCell(row, 3, tview.NewTableCell("Confirmed").SetAlign(tview.AlignCenter))
		} else {
			historyTable.SetCell(row, 3, tview.NewTableCell("Unconfirmed").SetAlign(tview.AlignCenter))
		}
		historyTable.SetCell(row, 1, tview.NewTableCell(tx.Direction.String()).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 2, tview.NewTableCell(tx.Amount).SetAlign(tview.AlignCenter))
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

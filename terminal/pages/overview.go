package pages

import (
	"context"
	"fmt"
	"time"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func overviewPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	overviewPage := tview.NewFlex().SetDirection(tview.FlexRow)

	var views []tview.Primitive
	var viewBoxes []*tview.Box

	errorTextView := primitives.WordWrappedTextView("")
	errorTextView.SetTextColor(helpers.DecredOrangeColor)

	displayError := func(errorMessage string) {
		overviewPage.RemoveItem(errorTextView)
		errorTextView.SetText(errorMessage)
		overviewPage.AddItem(errorTextView, 2, 0, false)
	}

	balanceViews, balanceViewBoxes := renderBalanceSection(overviewPage, wallet, displayError)
	views = append(views, balanceViews...)
	viewBoxes = append(viewBoxes, balanceViewBoxes...)

	overviewPage.AddItem(nil, 1, 0, false) // em

	recentActivityViews, recentActivityViewBoxes := renderRecentActivity(overviewPage, wallet, displayError)
	views = append(views, recentActivityViews...)
	viewBoxes = append(viewBoxes, recentActivityViewBoxes...)

	hintTextView.SetText("TIP: Scroll recent activity table with ARROW KEYS. Return to navigation menu with ESC")

	tabAndEscInputCapture := func(viewIndex int) func(event *tcell.EventKey) *tcell.EventKey {
		var nextView, previousView tview.Primitive
		if viewIndex == len(views)-1 {
			// this is last view, next view would be the first view
			nextView = views[0]
		} else {
			nextView = views[viewIndex+1]
		}
		if viewIndex == 0 {
			// this is first view, previous view will be last view
			previousView = views[len(views)-1]
		} else {
			previousView = views[viewIndex-1]
		}

		return func(event *tcell.EventKey) (nextEvent *tcell.EventKey) {
			if event.Key() == tcell.KeyEsc {
				clearFocus()
			} else if event.Key() == tcell.KeyTab {
				setFocus(nextView)
			} else if event.Key() == tcell.KeyBacktab {
				setFocus(previousView)
			} else {
				nextEvent = event
			}
			return
		}
	}

	for i, viewBox := range viewBoxes {
		viewBox.SetInputCapture(tabAndEscInputCapture(i))
	}

	setFocus(overviewPage)

	return overviewPage
}

func renderBalanceSection(overviewPage *tview.Flex, wallet walletcore.Wallet, displayError func(message string)) (views []tview.Primitive, viewBoxes []*tview.Box) {
	balanceTitleTextView := primitives.NewLeftAlignedTextView("Balance")
	overviewPage.AddItem(balanceTitleTextView, 2, 0, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		displayError(err.Error())
		return
	}

	overviewPage.AddItem(primitives.NewLeftAlignedTextView(walletcore.WalletBalance(accounts)), 2, 0, false)

	return
}

func renderRecentActivity(overviewPage *tview.Flex, wallet walletcore.Wallet, displayError func(message string)) (views []tview.Primitive, viewBoxes []*tview.Box) {
	overviewPage.AddItem(primitives.NewLeftAlignedTextView("-Recent Activity-").SetTextColor(helpers.DecredLightBlueColor), 1, 0, false)

	txns, _, err := wallet.TransactionHistory(context.Background(), -1, 5)
	if err != nil {
		displayError(err.Error())
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
	currentDate := time.Now().In(loc).Add(1 * time.Hour)
	timeDifference, _ := time.ParseDuration("24h")

	var confirmations int32 
	confirmations = walletcore.DefaultRequiredConfirmations
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
		transactionDate := time.Unix(tx.Timestamp, 0).In(loc).Add(1 * time.Hour)
		transactionDuration := currentDate.Sub(transactionDate)
	   	dateOutput  := strings.Split(tx.FormattedTime, " ")

	    if transactionDuration > timeDifference {
	    	historyTable.SetCell(row, 0, tview.NewTableCell(fmt.Sprintln(dateOutput[0])).SetAlign(tview.AlignCenter))
		}else{
	    	historyTable.SetCell(row, 0, tview.NewTableCell(fmt.Sprintln(dateOutput[1])).SetAlign(tview.AlignCenter))
		}

		if txns.Confirmations > confirmations{
			historyTable.SetCell(row, 3, tview.NewTableCell("Confirmed").SetAlign(tview.AlignCenter))
		}else{
			historyTable.SetCell(row, 3, tview.NewTableCell("Unconfirmed").SetAlign(tview.AlignCenter))
		}
		historyTable.SetCell(row, 1, tview.NewTableCell(tx.Direction.String()).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 2, tview.NewTableCell(tx.Amount).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 4, tview.NewTableCell(tx.Type).SetAlign(tview.AlignCenter))
	}

	overviewPage.AddItem(historyTable, 0, 1, true)
	views = append(views, historyTable)
	viewBoxes = append(viewBoxes, historyTable.Box)

	return
}

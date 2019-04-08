package pages

import (
	"context"
	"fmt"
	"time"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func overviewPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	overviewPage := tview.NewFlex().SetDirection(tview.FlexRow)

	var views []tview.Primitive
	var viewBoxes []*tview.Box

	balanceViews, balanceViewBoxes := renderBalanceSection(overviewPage, wallet)
	views = append(views, balanceViews...)
	viewBoxes = append(viewBoxes, balanceViewBoxes...)

	overviewPage.AddItem(nil, 1, 0, false) // em

	recentActivityViews, recentActivityViewBoxes := renderRecentActivity(overviewPage, wallet)
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

func renderBalanceSection(overviewPage *tview.Flex, wallet walletcore.Wallet) (views []tview.Primitive, viewBoxes []*tview.Box) {
	balanceTitleTextView := primitives.NewLeftAlignedTextView("Balance")
	overviewPage.AddItem(balanceTitleTextView, 2, 0, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		overviewPage.AddItem(primitives.NewCenterAlignedTextView(err.Error()), 3, 0, false)
		return
	}

	balanceTable := primitives.NewTable()
	balanceTable.SetBorders(false).SetFixed(1, 0)

	tableHeight := len(accounts) + 1 // 1 row for each account, plus the table header
	overviewPage.AddItem(balanceTable, tableHeight, 0, true)
	views = append(views, balanceTable)
	viewBoxes = append(viewBoxes, balanceTable.Box)

	toggleBalanceDisplay := func() {
		if len(accounts) == 1 {
			displaySingleAccountSimpleBalance(accounts[0], balanceTable)
		} else {
			displayMultipleAccountsSimpleBalance(accounts, balanceTable)
		}
	}

	toggleBalanceDisplay()

	return
}

func displaySingleAccountSimpleBalance(account *walletcore.Account, balanceTable *primitives.Table) {
	if account.Balance.Total == account.Balance.Spendable {
		// show only total since it is equal to spendable
		balanceTable.SetCellSimple(0, 0, walletcore.NormalizeBalance(account.Balance.Total.ToCoin()))
	} else {
		balanceTable.SetCellSimple(0, 0, "Total")
		balanceTable.SetCellRightAlign(0, 1, walletcore.NormalizeBalance(account.Balance.Total.ToCoin()))
		balanceTable.SetCellSimple(1, 0, "Spendable")
		balanceTable.SetCellRightAlign(1, 1, walletcore.NormalizeBalance(account.Balance.Spendable.ToCoin()))
	}
}

func displayMultipleAccountsSimpleBalance(accounts []*walletcore.Account, balanceTable *primitives.Table) {
	// draw table header
	balanceTable.SetHeaderCell(0, 0, "Account Name").
		SetHeaderCell(0, 1, "Balance").
		SetHeaderCell(0, 2, "Spendable")

	for i, account := range accounts {
		row := i + 1
		balanceTable.SetCellCenterAlign(row, 0, account.Name).
			SetCellRightAlign(row, 1, walletcore.NormalizeBalance(account.Balance.Total.ToCoin())).
			SetCellRightAlign(row, 2, walletcore.NormalizeBalance(account.Balance.Spendable.ToCoin()))
	}
}

func renderRecentActivity(overviewPage *tview.Flex, wallet walletcore.Wallet) (views []tview.Primitive, viewBoxes []*tview.Box) {
	overviewPage.AddItem(primitives.NewLeftAlignedTextView("Recent Activity"), 1, 0, false)

	txns, _, err := wallet.TransactionHistory(context.Background(), -1, 5)
	if err != nil {
		overviewPage.AddItem(primitives.NewCenterAlignedTextView(err.Error()), 3, 0, false)
		return
	}

	historyTable := primitives.NewTable()
	historyTable.SetBorders(false).SetFixed(1, 0)

	// historyTable header
	historyTable.SetHeaderCell(0, 0, "#")
	historyTable.SetHeaderCell(0, 1, "Date")
	historyTable.SetHeaderCell(0, 3, "Direction")
	historyTable.SetHeaderCell(0, 2, "Amount")
	historyTable.SetHeaderCell(0, 4, "Status")
	historyTable.SetHeaderCell(0, 5, "Type")

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
			fmt.Println(err.Error())
			return
		}
		transactionDate := time.Unix(tx.Timestamp, 0).In(loc).Add(1 * time.Hour)
		transactionDuration := currentDate.Sub(transactionDate)
	   	date := strings.Split(tx.FormattedTime, " ")

	    if transactionDuration > timeDifference {
	    	historyTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintln(date[0], date[2])).SetAlign(tview.AlignCenter))
		}else{
	    	historyTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintln(date[1], date[2])).SetAlign(tview.AlignCenter))
		}
		if txns.Confirmations > confirmations{
			historyTable.SetCell(row, 4, tview.NewTableCell("Confirmed").SetAlign(tview.AlignCenter))
		}else{
			historyTable.SetCell(row, 4, tview.NewTableCell("Unconfirmed").SetAlign(tview.AlignCenter))
		}
		historyTable.SetCellSimple(row, 0, fmt.Sprintf("%d.", row))
		historyTable.SetCell(row, 3, tview.NewTableCell(tx.Direction.String()).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 2, tview.NewTableCell(tx.Amount).SetAlign(tview.AlignRight))
		historyTable.SetCell(row, 5, tview.NewTableCell(tx.Type).SetAlign(tview.AlignCenter))
	}

	overviewPage.AddItem(historyTable, 0, 1, true)
	views = append(views, historyTable)
	viewBoxes = append(viewBoxes, historyTable.Box)

	return
}

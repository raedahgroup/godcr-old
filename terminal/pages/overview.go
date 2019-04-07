package pages

import (
	"context"
	"fmt"

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

	hintTextView.SetText("TIP: Move around with TAB and SHIFT+TAB. Scroll tables with ARROW KEYS. " +
		"Return to navigation menu with Esc")

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

	var toggleBalanceForm *tview.Form
	var showSimpleBalanceNext bool

	toggleBalanceDisplay := func() {
		balanceTable.Clear()

		if showSimpleBalanceNext {
			showSimpleBalanceNext = false
			balanceTitleTextView.SetText("Balance")
			toggleBalanceForm.GetButton(0).SetLabel("Show Detailed Balance")

			if len(accounts) == 1 {
				displaySingleAccountSimpleBalance(accounts[0], balanceTable)
			} else {
				displayMultipleAccountsSimpleBalance(accounts, balanceTable)
			}
		} else {
			showSimpleBalanceNext = true
			balanceTitleTextView.SetText("Balance (Detailed)")
			toggleBalanceForm.GetButton(0).SetLabel("Show Simple Balance")
			displayDetailedAccountsBalances(accounts, balanceTable)
		}
	}

	// display button to toggle balance display, embed button in form so it doesn't fill screen width
	toggleBalanceForm = tview.NewForm().AddButton("Show Detailed Balance", toggleBalanceDisplay)
	toggleBalanceForm.SetBorderPadding(0, 0, 0, 0)
	toggleBalanceForm.SetItemPadding(0)

	overviewPage.AddItem(toggleBalanceForm, 2, 0, false)
	views = append(views, toggleBalanceForm.GetButton(0))
	viewBoxes = append(viewBoxes, toggleBalanceForm.GetButton(0).Box)

	showSimpleBalanceNext = true
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

func displayDetailedAccountsBalances(accounts []*walletcore.Account, balanceTable *primitives.Table) {
	// draw table header
	balanceTable.SetHeaderCell(0, 0, "Account Name").
		SetHeaderCell(0, 1, "Balance").
		SetHeaderCell(0, 2, "Spendable").
		SetHeaderCell(0, 3, "Locked").
		SetHeaderCell(0, 4, "Voting Authority").
		SetHeaderCell(0, 5, "Unconfirmed")

	for i, account := range accounts {
		row := i + 1
		balanceTable.SetCellCenterAlign(row, 0, account.Name).
			SetCellRightAlign(row, 1, walletcore.NormalizeBalance(account.Balance.Total.ToCoin())).
			SetCellRightAlign(row, 2, walletcore.NormalizeBalance(account.Balance.Spendable.ToCoin())).
			SetCellCenterAlign(row, 3, account.Balance.LockedByTickets.String()).
			SetCellCenterAlign(row, 4, account.Balance.VotingAuthority.String()).
			SetCellCenterAlign(row, 5, account.Balance.Unconfirmed.String())
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
	historyTable.SetHeaderCell(0, 4, "Direction")
	historyTable.SetHeaderCell(0, 2, "Amount")
	historyTable.SetHeaderCell(0, 3, "Fee")
	historyTable.SetHeaderCell(0, 5, "Type")
	historyTable.SetHeaderCell(0, 6, "Hash")

	for _, tx := range txns {
		row := historyTable.GetRowCount()
		if row >= 5 {
			break
		}

		historyTable.SetCellSimple(row, 0, fmt.Sprintf("%d.", row))
		historyTable.SetCell(row, 1, tview.NewTableCell(tx.FormattedTime).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 4, tview.NewTableCell(tx.Direction.String()).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 2, tview.NewTableCell(tx.Amount).SetAlign(tview.AlignRight))
		historyTable.SetCell(row, 3, tview.NewTableCell(tx.Fee).SetAlign(tview.AlignRight))
		historyTable.SetCell(row, 5, tview.NewTableCell(tx.Type).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 6, tview.NewTableCell(tx.Hash).SetAlign(tview.AlignCenter))
	}

	overviewPage.AddItem(historyTable, 0, 1, true)
	views = append(views, historyTable)
	viewBoxes = append(viewBoxes, historyTable.Box)

	return
}

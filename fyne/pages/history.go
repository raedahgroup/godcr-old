package pages

import (
	"context"
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type historyPageLoader struct {
	ctx                    context.Context
	wallet                 walletcore.Wallet
	updatePageOnMainWindow func()
	transactionTable       *widgets.Table
	loadTransactionsError  string
	transactionCount       int
}

// Load initializes the page views and updates the app window before and/or after loading data
func (page *historyPageLoader) Load(ctx context.Context, wallet walletcore.Wallet, updatePageOnMainWindow func(object fyne.CanvasObject)) {
	page.wallet = wallet
	page.updatePageOnMainWindow = page.makePageUpdateFunc(updatePageOnMainWindow)
	page.transactionTable = widgets.NewTable()

	// update window now, then fetch recent activity in background, as it "may" take some time
	page.updatePageOnMainWindow()
	go page.loadTransactions()
}

func (page *historyPageLoader) loadTransactions() {
	// update main window after fetching recent activity
	defer page.updatePageOnMainWindow()

	page.transactionTable.Clear()
	var err error
	page.transactionCount, err = page.wallet.TransactionCount(nil)
	if err != nil {
		page.loadTransactionsError = err.Error()
		return
	}
	txns, err := page.wallet.TransactionHistory(-1, int32(page.transactionCount), nil)
	if err != nil {
		page.loadTransactionsError = err.Error()
		return
	}
	leftAlign := []int{3}
	page.transactionTable.AddRowHeader("Account", "Date", "Type", "Direction", "Amount", "Fee", "Status", "Hash")
	for _, tx := range txns {
		trimmedHash := tx.Hash[:len(tx.Hash)/2] + "..."
		page.transactionTable.AddRowWithButtonSupport(tx.Hash,
			7,
			leftAlign,
			tx.AccountName(),
			tx.LongTime,
			tx.Type,
			tx.Direction.String(),
			dcrutil.Amount(tx.Amount).String(),
			dcrutil.Amount(tx.Fee).String(),
			tx.Status,
			trimmedHash,
		)

	}

}

// makePageUpdateFunc creates a wrapper function around `updatePageOnMainWindow`
// to update the app window when relevant changes are made to the page content
func (page *historyPageLoader) makePageUpdateFunc(updatePageOnMainWindow func(object fyne.CanvasObject)) func() {
	return func() {
		var recentActivityObject fyne.CanvasObject
		if page.loadTransactionsError != "" {
			recentActivityObject = widget.NewLabel(fmt.Sprintf("Error History: %s", page.loadTransactionsError))
		} else {
			recentActivityObject = page.transactionTable.CondensedTable()
		}
		pageViews := widget.NewVBox(
			widgets.NewVSpacer(20),

			recentActivityObject,
		)
		updatePageOnMainWindow(pageViews)
	}
}

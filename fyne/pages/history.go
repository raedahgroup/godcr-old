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
	ctx                      context.Context
	wallet                   walletcore.Wallet
	updatePageOnMainWindow   func()
	recentActivityTable      *widgets.Table
	fetchRecentActivityError string
}

// Load initializes the page views and updates the app window before and/or after loading data
func (page *historyPageLoader) Load(ctx context.Context, wallet walletcore.Wallet, updatePageOnMainWindow func(object fyne.CanvasObject)) {
	page.wallet = wallet
	page.updatePageOnMainWindow = page.makePageUpdateFunc(updatePageOnMainWindow)
	page.recentActivityTable = widgets.NewTable()

	// update window now, then fetch recent activity in background, as it "may" take some time
	page.updatePageOnMainWindow()
	go page.fetchRecentActivity()
}

func (page *historyPageLoader) fetchRecentActivity() {
	// update main window after fetching recent activity
	defer page.updatePageOnMainWindow()

	page.recentActivityTable.Clear()
	no, err := page.wallet.TransactionCount(nil)
	if err != nil {
		page.fetchRecentActivityError = err.Error()
		return
	}
	txns, err := page.wallet.TransactionHistory(-1, int32(no), nil)
	if err != nil {
		page.fetchRecentActivityError = err.Error()
		return
	}

	page.recentActivityTable.AddRowSimple("#", "Date", "Direction", "Amount", "Fee", "Type", "Hash")
	for i, tx := range txns {
		trimmedHash := tx.Hash[:len(tx.Hash)/2] + "..."
		page.recentActivityTable.AddRowWithHashClick(
			tx.Hash,
			fmt.Sprintf("%d", i+1),
			tx.LongTime,
			tx.Direction.String(),
			dcrutil.Amount(tx.Amount).String(),
			dcrutil.Amount(tx.Fee).String(),
			tx.Type,
			trimmedHash,
		)
	}
}

// makePageUpdateFunc creates a wrapper function around `updatePageOnMainWindow`
// to update the app window when relevant changes are made to the page content
func (page *historyPageLoader) makePageUpdateFunc(updatePageOnMainWindow func(object fyne.CanvasObject)) func() {
	return func() {
		var recentActivityObject fyne.CanvasObject
		if page.fetchRecentActivityError != "" {
			recentActivityObject = widget.NewLabel(fmt.Sprintf("Error History: %s", page.fetchRecentActivityError))
		} else {
			recentActivityObject = page.recentActivityTable.CondensedTable()
		}

		pageViews := widget.NewVBox(
			widgets.NewVSpacer(20),
			recentActivityObject,
		)
		updatePageOnMainWindow(pageViews)
	}
}

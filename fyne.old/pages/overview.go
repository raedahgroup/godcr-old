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

type overviewPageLoader struct {
	ctx                    context.Context
	wallet                 walletcore.Wallet
	updatePageOnMainWindow func()

	balanceSectionTitle *widget.Label
	balanceLabel        *widget.Label
	fetchBalanceError   string

	recentActivitySectionTitle *widget.Label
	recentActivityTable        *widgets.Table
	fetchRecentActivityError   string
}

// Load initializes the page views and updates the app window before and/or after loading data
func (page *overviewPageLoader) Load(ctx context.Context, wallet walletcore.Wallet, updatePageOnMainWindow func(object fyne.CanvasObject)) {
	page.wallet = wallet
	page.updatePageOnMainWindow = page.makePageUpdateFunc(updatePageOnMainWindow)

	// init balance views
	page.balanceSectionTitle = widget.NewLabelWithStyle("Current Total Balance", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	page.fetchBalance()

	// init recent activity views
	page.recentActivitySectionTitle = widget.NewLabelWithStyle("Recent Activity", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	page.recentActivityTable = widgets.NewTable()

	// update window now, then fetch recent activity in background, as it "may" take some time
	page.updatePageOnMainWindow()
	go page.fetchRecentActivity()
}

func (page *overviewPageLoader) fetchBalance() {
	accounts, err := page.wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		page.fetchBalanceError = err.Error()
		return
	}

	var totalBalance dcrutil.Amount
	for _, account := range accounts {
		totalBalance += account.Balance.Total
	}

	page.balanceLabel = widget.NewLabel(walletcore.NormalizeBalance(accounts[0].Balance.Total.ToCoin()))
}

func (page *overviewPageLoader) fetchRecentActivity() {
	// update main window after fetching recent activity
	defer page.updatePageOnMainWindow()

	page.recentActivityTable.Clear()

	txns, err := page.wallet.TransactionHistory(0, 5, nil)
	if err != nil {
		page.fetchRecentActivityError = err.Error()
		return
	}

	page.recentActivityTable.AddRowSimple("#", "Date", "Direction", "Amount", "Fee", "Type", "Hash")
	for i, tx := range txns {
		page.recentActivityTable.AddRowSimple(
			fmt.Sprintf("%d", i+1),
			tx.ShortTime,
			tx.Direction.String(),
			dcrutil.Amount(tx.Amount).String(),
			dcrutil.Amount(tx.Fee).String(),
			tx.Type,
			tx.Hash,
		)
	}
}

// makePageUpdateFunc creates a wrapper function around `updatePageOnMainWindow`
// to update the app window when relevant changes are made to the page content
func (page *overviewPageLoader) makePageUpdateFunc(updatePageOnMainWindow func(object fyne.CanvasObject)) func() {
	return func() {
		var balanceView fyne.CanvasObject
		if page.fetchBalanceError != "" {
			balanceView = widget.NewLabel(fmt.Sprintf("Error fetching balance: %s", page.fetchBalanceError))
		} else {
			balanceView = page.balanceLabel
		}

		var recentActivityObject fyne.CanvasObject
		if page.fetchBalanceError != "" {
			recentActivityObject = widget.NewLabel(fmt.Sprintf("Error fetching recent activity: %s", page.fetchRecentActivityError))
		} else {
			recentActivityObject = page.recentActivityTable.CondensedTable()
		}

		pageViews := widget.NewVBox(
			page.balanceSectionTitle,
			balanceView,
			widgets.NewVSpacer(20),
			page.recentActivitySectionTitle,
			recentActivityObject,
		)
		updatePageOnMainWindow(pageViews)
	}
}

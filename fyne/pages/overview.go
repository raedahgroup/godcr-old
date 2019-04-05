package pages

import (
	"context"
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type overviewPageLoader struct {
	ctx                    context.Context
	wallet                 walletcore.Wallet
	updatePageOnMainWindow func()

	balanceSectionTitle *widget.Label
	showDetailsCheckbox *widget.Check
	balanceTable        *widgets.Table
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
	page.balanceSectionTitle = widget.NewLabelWithStyle("- Balance -", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	page.balanceTable = widgets.NewTable()
	page.showDetailsCheckbox = widget.NewCheck("Show details", func(showDetails bool) {
		page.fetchBalance(showDetails)
	})

	// init recent activity views
	page.recentActivitySectionTitle = widget.NewLabelWithStyle("- Recent Activity -", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	page.recentActivityTable = widgets.NewTable()

	// fetch balance before updating window, shouldn't be a long-running operation
	// the page.fetchBalance() method will update the window when done
	page.fetchBalance(false)

	// fetch recent activity in background, "may" take some time
	go page.fetchRecentActivity()
}

func (page *overviewPageLoader) fetchBalance(detailed bool) {
	// update main window after fetching balance
	defer page.updatePageOnMainWindow()

	page.balanceTable.Clear()

	accounts, err := page.wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		page.fetchBalanceError = err.Error()
		return
	}

	if len(accounts) == 1 && !detailed {
		account := accounts[0]
		if account.Balance.Total == account.Balance.Spendable {
			// show only total since it is equal to spendable
			page.balanceTable.AddRowSimple(walletcore.NormalizeBalance(account.Balance.Total.ToCoin()))
		} else {
			page.balanceTable.AddRowSimple("Total", walletcore.NormalizeBalance(account.Balance.Total.ToCoin()))
			page.balanceTable.AddRowSimple("Spendable", walletcore.NormalizeBalance(account.Balance.Spendable.ToCoin()))
		}
		return
	}

	// if there are more than 1 account or it's 1 account but we're required to show details,
	// let's use a proper table with headers
	columnHeaders := []string{
		"Account",
		"Total",
		"Spendable",
	}
	if detailed {
		columnHeaders = append(columnHeaders, "Locked")
		columnHeaders = append(columnHeaders, "Voting Authority")
		columnHeaders = append(columnHeaders, "Unconfirmed")
	}
	page.balanceTable.AddRowSimple(columnHeaders...)

	for _, account := range accounts {
		rowValues := []string{
			account.Name,
			walletcore.NormalizeBalance(account.Balance.Total.ToCoin()),
			walletcore.NormalizeBalance(account.Balance.Spendable.ToCoin()),
		}
		if detailed {
			rowValues = append(rowValues, account.Balance.LockedByTickets.String())
			rowValues = append(rowValues, account.Balance.VotingAuthority.String())
			rowValues = append(rowValues, account.Balance.Unconfirmed.String())
		}
		page.balanceTable.AddRowSimple(rowValues...)
	}
}

func (page *overviewPageLoader) fetchRecentActivity() {
	// update main window after fetching recent activity
	defer page.updatePageOnMainWindow()

	page.recentActivityTable.Clear()

	txns, _, err := page.wallet.TransactionHistory(context.Background(), -1, 5)
	if err != nil {
		page.fetchRecentActivityError = err.Error()
		return
	}

	page.recentActivityTable.AddRowSimple("#", "Date", "Direction", "Amount", "Fee", "Type", "Hash")
	for i, tx := range txns {
		page.recentActivityTable.AddRowSimple(
			fmt.Sprintf("%d", i+1),
			tx.FormattedTime,
			tx.Direction.String(),
			tx.Amount,
			tx.Fee,
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
			balanceView = page.balanceTable.CondensedTable()
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
			page.showDetailsCheckbox,
			widgets.NewVSpacer(20),
			page.recentActivitySectionTitle,
			recentActivityObject,
		)
		updatePageOnMainWindow(pageViews)
	}
}

package pages

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func overviewPageContent(wallet walletcore.Wallet, updatePageOnMainWindow func(object fyne.CanvasObject)) {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		errorLabel := widget.NewLabel(fmt.Sprintf("Error loading accounts overview: %s", err.Error()))
		updatePageOnMainWindow(errorLabel)
	}

	var showDetailsCheckbox *widget.Check
	balanceTable := widgets.NewTable()

	resizePage := func() {
		overviewPageContent := widget.NewVBox(showDetailsCheckbox, balanceTable.CondensedTable())
		updatePageOnMainWindow(overviewPageContent)
	}

	showDetailsCheckbox = widget.NewCheck("Show Detailed Balance", func(showDetails bool) {
		balanceTable.Clear()
		displayBalance(accounts, balanceTable, showDetails, resizePage)
	})

	overviewPageContent := widget.NewVBox(showDetailsCheckbox, balanceTable.CondensedTable())
	updatePageOnMainWindow(overviewPageContent)

	displayBalance(accounts, balanceTable, false, resizePage)
}

func displayBalance(accounts []*walletcore.Account, balanceTable *widgets.Table, detailed bool, resizePage func()) {
	// resize balance table when done
	defer resizePage()

	if len(accounts) == 1 && !detailed {
		account := accounts[0]
		if account.Balance.Total == account.Balance.Spendable {
			// show only total since it is equal to spendable
			balanceTable.AddRowSimple(walletcore.NormalizeBalance(account.Balance.Total.ToCoin()))
		} else {
			balanceTable.AddRowSimple("Total", walletcore.NormalizeBalance(account.Balance.Total.ToCoin()))
			balanceTable.AddRowSimple("Spendable", walletcore.NormalizeBalance(account.Balance.Spendable.ToCoin()))
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
	balanceTable.AddRowSimple(columnHeaders...)

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
		balanceTable.AddRowSimple(rowValues...)
	}
}

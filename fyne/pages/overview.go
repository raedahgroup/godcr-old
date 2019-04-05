package pages

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func overviewPageContent(wallet walletcore.Wallet) fyne.CanvasObject {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return widget.NewLabel(fmt.Sprintf("Error loading accounts overview: %s", err.Error()))
	}

	balanceTable := widgets.NewTable()

	showDetailsCheckbox := widget.NewCheck("Show Detailed Balance", func(showDetails bool) {
		displayBalance(accounts, balanceTable, showDetails)
	})

	displayBalance(accounts, balanceTable, false)

	return widget.NewVBox(showDetailsCheckbox, balanceTable)
}

func displayBalance(accounts []*walletcore.Account, balanceTable *widgets.Table, detailed bool) {
	if len(accounts) == 0 && !detailed {
		account := accounts[0]
		if account.Balance.Total == account.Balance.Spendable {
			// show only total since it is equal to spendable
			balanceTable.AddObject(widget.NewLabel(walletcore.NormalizeBalance(account.Balance.Total.ToCoin())))
		} else {
			balanceTable.AddRowSimple("Total", walletcore.NormalizeBalance(account.Balance.Total.ToCoin()))
			balanceTable.AddRowSimple("Spendable", walletcore.NormalizeBalance(account.Balance.Spendable.ToCoin()))
		}
		return
	}

	// if there are more than 1 account or it's 1 account but we're required to show details,
	// let's use a proper table with headers
	columnHeaders := []string{
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
			walletcore.NormalizeBalance(account.Balance.Total.ToCoin()),
			walletcore.NormalizeBalance(account.Balance.Spendable.ToCoin()),
		}
		if detailed {
			rowValues = append(rowValues, account.Balance.LockedByTickets.String())
			rowValues = append(rowValues, account.Balance.VotingAuthority.String())
			rowValues = append(rowValues, account.Balance.Unconfirmed.String())
		}
	}
}

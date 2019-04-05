package pages

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/layout"

	"github.com/raedahgroup/godcr/app/walletcore"
)

func overviewPageContent(wallet walletcore.Wallet) fyne.CanvasObject {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return widget.NewLabel(fmt.Sprintf("Error loading accounts overview: %s", err.Error()))
	}

	// use grid for balance table
	balanceTable := fyne.NewContainerWithLayout(layout.NewGridLayout(1))
	displayBalance(accounts, balanceTable, false)

	return balanceTable
}

func displayBalance(accounts []*walletcore.Account, balanceTable *fyne.Container, detailed bool) {
	if len(accounts) == 0 && !detailed {
		account := accounts[0]
		if account.Balance.Total == account.Balance.Spendable {
			// show only total since it is equal to spendable
			balanceTable.AddObject(widget.NewLabel(walletcore.NormalizeBalance(account.Balance.Total.ToCoin())))
		} else {
			balanceTable.SetCellSimple(0, 0, "Total")
			balanceTable.SetCellRightAlign(0, 1, walletcore.NormalizeBalance(account.Balance.Total.ToCoin()))
			balanceTable.SetCellSimple(1, 0, "Spendable")
			balanceTable.SetCellRightAlign(1, 1, walletcore.NormalizeBalance(account.Balance.Spendable.ToCoin()))
		}
	}
}

// Package pages contains implementations of the various pages that
// can be viewed on the desktop GUI when using the "fyne" library.
//
// This file contains the implementation of the accounts overview page.
package pages

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/log"
)

func accountsPage(wallet app.WalletMiddleware) fyne.CanvasObject {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		log.Error(err.Error())
		return widget.NewLabel(fmt.Sprintf("Unable to retrieve accounts information: %s", err))
	}
	accountsGroup := widget.NewGroup("Accounts")
	for _, account := range accounts {
		accountName := widget.NewLabel(account.Name)
		accountBalance := widget.NewLabel(account.Balance.String())
		accountsGroup.Append(widget.NewVBox(widget.NewHBox(accountName, accountBalance)))
	}
	return accountsGroup
}

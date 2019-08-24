// Package pages contains implementations of the various pages that
// can be viewed on the desktop GUI when using the "fyne" library.
//
// This file contains the implementation of the accounts overview page.
package pages

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/log"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

const defaultItemsSpacing = 10

func accountsPage(wallet app.WalletMiddleware) fyne.CanvasObject {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		log.Error(err.Error())
		return widget.NewLabel(fmt.Sprintf("Unable to retrieve accounts information: %s", err))
	}
	netType := wallet.NetType()
	pageTitle := widget.NewLabelWithStyle("Accounts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	accountsBox := widget.NewVBox(pageTitle)
	for _, account := range accounts {
		accountsBox.Append(widgets.NewVSpacer(defaultItemsSpacing))
		hSpace := widgets.NewHSpacer(defaultItemsSpacing)
		detailSection := widget.NewHBox(hSpace, createDetailSection(account, netType))
		toggleButton := createToggleButton(detailSection)
		overviewSection := widget.NewHBox(toggleButton, createOverviewSection(account))
		accountsBox.Append(widget.NewVBox(overviewSection, detailSection))
	}
	addAccountForm := createAddAccountForm(wallet)
	addAccountButton := createAddAccountButton(addAccountForm)
	pageLeftMargin := widgets.NewHSpacer(defaultItemsSpacing)
	return widget.NewHBox(
		pageLeftMargin,
		accountsBox,
		pageLeftMargin,
		widget.NewVBox(widget.NewHBox(addAccountButton, layout.NewSpacer()), addAccountForm),
	)
}

func toggleCanvasObjectVisibility(target fyne.CanvasObject) {
	if target.Visible() {
		target.Hide()
	} else {
		target.Show()
	}
}

func createAddAccountButton(form fyne.CanvasObject) fyne.CanvasObject {
	return &widget.Button{
		Text:  "Add account",
		Style: widget.PrimaryButton,
		OnTapped: func() {
			toggleCanvasObjectVisibility(form)
		},
	}
}

func createAddAccountForm(wallet walletcore.Wallet) fyne.CanvasObject {
	accountNameEntry := &widget.Entry{
		PlaceHolder: "Account name",
	}
	accountName := &widget.FormItem{
		Text:   "Account name",
		Widget: accountNameEntry,
	}
	walletPassphraseEntry := &widget.Entry{
		PlaceHolder: "Wallet passphrase",
		Password:    true,
	}
	walletPassphrase := &widget.FormItem{
		Text:   "Wallet passphrase",
		Widget: walletPassphraseEntry,
	}
	errorDisplay := widget.NewLabel("")
	form := &widget.Form{
		Items: []*widget.FormItem{accountName, walletPassphrase},
	}
	resetForm := func() {
		accountNameEntry.SetText("")
		walletPassphraseEntry.SetText("")
		errorDisplay.SetText("")
		errorDisplay.Hide()
	}
	form.OnCancel = resetForm
	form.OnSubmit = func() {
		_, err := wallet.NextAccount(accountNameEntry.Text, walletPassphraseEntry.Text)
		if err != nil {
			errorDisplay.SetText(fmt.Sprintf("Failed to create account. Reason: %s", err))
			errorDisplay.Show()
		} else {
			resetForm()
		}
	}
	return widget.NewVBox(form, errorDisplay)
}

func createOverviewSection(account *walletcore.Account) fyne.CanvasObject {
	accountName := widget.NewLabel(account.Name)
	totalBalance := widget.NewLabel(fmt.Sprintf("Total balance: %s", account.Balance.Total))
	spendableBalance := widget.NewLabel(fmt.Sprintf("Spendable balance: %s", account.Balance.Spendable))
	return widget.NewHBox(accountName, layout.NewSpacer(), totalBalance, spendableBalance)
}

func createDetailSection(account *walletcore.Account, netType string) fyne.CanvasObject {
	properties := widget.NewHBox(widget.NewLabelWithStyle("Properties", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	hSpace := widgets.NewHSpacer(defaultItemsSpacing)
	accountNumber := widget.NewHBox(
		widget.NewLabel("Account number"),
		layout.NewSpacer(),
		widget.NewLabel(fmt.Sprintf("%d", account.Number)))
	hdPath := walletcore.MainnetHDPath
	if netType == "testnet3" {
		hdPath = walletcore.TestnetHDPath
	}
	hdPathBox := widget.NewHBox(widget.NewLabel("HD Path"), layout.NewSpacer(), widget.NewLabel(hdPath))
	keysText := fmt.Sprintf(
		"%d external, %d internal, %d imported",
		account.ExternalKeyCount, account.InternalKeyCount, account.ImportedKeyCount)
	keys := widget.NewHBox(widget.NewLabel("Keys"), layout.NewSpacer(), widget.NewLabel(keysText))
	detailsBox := widget.NewVBox(properties, widget.NewHBox(hSpace, widget.NewVBox(accountNumber, hdPathBox, keys)))
	return detailsBox
}

func createToggleButton(target fyne.CanvasObject) *widget.Button {
	button := widget.Button{}
	setButtonStyle := func() {
		if target.Visible() {
			button.SetText("-")
			button.Style = widget.PrimaryButton
		} else {
			button.SetText("+")
			button.Style = widget.DefaultButton
		}
	}
	setButtonStyle()
	button.OnTapped = func() {
		toggleCanvasObjectVisibility(target)
		setButtonStyle()
	}
	return &button
}

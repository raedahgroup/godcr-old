// Package pages contains implementations of the various pages that
// can be viewed on the desktop GUI when using the "fyne" library.
//
// This file contains the implementation of the accounts overview page.
package pages

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
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
	pageTitle := widget.NewLabelWithStyle("Accounts", fyne.TextAlignLeading, fyne.TextStyle{Bold:true, Italic:true})
	pageBox := widget.NewVBox(pageTitle)
	for _, account := range accounts {
		pageBox.Append(widgets.NewVSpacer(defaultItemsSpacing))
		hSpace := widgets.NewHSpacer(defaultItemsSpacing)
		detailSection := widget.NewHBox(hSpace, createDetailSection(account, netType))
		toggleButton := createToggleButton(detailSection)
		overviewSection := widget.NewHBox(toggleButton, createOverviewSection(account))
		pageBox.Append(widget.NewVBox(overviewSection, detailSection))
	}
	pageLeftMargin := widgets.NewHSpacer(defaultItemsSpacing)
	return widget.NewHBox(pageLeftMargin, pageBox)
}

func createOverviewSection(account *walletcore.Account) fyne.CanvasObject {
	accountName := widget.NewLabel(account.Name)
	accountBalance := widget.NewLabelWithStyle(account.Balance.String(), fyne.TextAlignTrailing, fyne.TextStyle{})
	return widget.NewHBox(accountName, layout.NewSpacer(), accountBalance)
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
		if target.Visible() {
			target.Hide()
		} else {
			target.Show()
		}
		setButtonStyle()
		canvas.Refresh(&button)
	}
	return &button
}
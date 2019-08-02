package pages

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type sendPageData struct {
	accountSelect *widget.Select
}

//both send page and receive page update would be in a function
var send sendPageData

func sendPage(wallet godcrApp.WalletMiddleware, window fyne.Window) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Sending Decred", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	errorLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	accountLabel := widget.NewLabelWithStyle("From:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		errorLabel.SetText("Could not retrieve account information" + err.Error())
		//todo: log to file
		fmt.Println(err.Error())
		errorLabel.Show()
		widget.Refresh(errorLabel)
	}

	var options []string
	for _, account := range accounts {
		options = append(options, account.String())
	}

	var button *widget.Button
	receive.accountSelect = widget.NewSelect(options, func(s string) {
		if button.Disabled() == true {
			button.Enable()
			widget.Refresh(button)
		}
	})

	address := widget.NewEntry()
	address.SetPlaceHolder("Destination Address")

	var spendUnconfirmed bool
	spendUnconfirmedCheck := widget.NewCheck("Spend Unconfirmed", func(check bool) {
		spendUnconfirmed = check
	})

	infoIcon := widget.NewToolbar(widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
		var button *widget.Button
		var popUp *widget.PopUp
		button = widget.NewButton("Got it", func() {
			popUp.Hide()
		})

		header := widget.NewLabelWithStyle("Send DCR", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		info := widget.NewLabelWithStyle("Input the destination wallet address and the amount to send funds", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
		data := widget.NewVBox(header, widgets.NewVSpacer(10), info,
			widget.NewHBox(layout.NewSpacer(), button))
		popUp = widget.NewModalPopUp(data, window.Canvas())

	}))

	return widget.NewVBox(widget.NewHBox(widgets.NewHSpacer(10), label, infoIcon), widget.NewLabel("kdmmmmmmmmmmmmmmmmmcldnvjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjfv"), spendUnconfirmedCheck)
}

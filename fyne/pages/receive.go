package pages

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
	"github.com/skip2/go-qrcode"
)

type receivePageData struct {
	accountSelect *widget.Select
}

var receive receivePageData

func receiveUpdates(wallet godcrApp.WalletMiddleware) {
	accounts, _ := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)

	var name []string
	for _, account := range accounts {
		name = append(name, account.String())
	}
	receive.accountSelect.Options = name
	widget.Refresh(receive.accountSelect)
}

//todo: should we make concurrent checks if users add a new account?
func receivePage(wallet godcrApp.WalletMiddleware, window fyne.Window) fyne.CanvasObject {
	qrImage := canvas.NewImageFromResource(theme.InfoIcon())
	qrImage.SetMinSize(fyne.NewSize(300, 300))
	qrImage.Hide()

	label := widget.NewLabelWithStyle("Receiving Funds", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	info := widget.NewLabelWithStyle("Each time you request a payment, a new address is created to protect your privacy.", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Italic: true})
	accountLabel := widget.NewLabelWithStyle("Account:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	generatedAddress := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	errorLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	generatedAddress.Hide()
	errorLabel.Hide()

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
		options = append(options, account.Name)
	}

	var button *widget.Button

	receive.accountSelect = widget.NewSelect(options, func(s string) {
		if button.Disabled() == true {
			button.Enable()
			widget.Refresh(button)
		}
	})

	var addr string

	copy := widget.NewToolbar(widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
		clipboard := window.Clipboard()
		clipboard.SetContent(addr)
	}))
	copy.Hide()

	button = widget.NewButton("Generate Address", func() {
		name, err := wallet.AccountNumber(receive.accountSelect.Selected)
		if err != nil {
			errorLabel.SetText("error getting account name, " + err.Error())
			errorLabel.Show()
			widget.Refresh(errorLabel)
			return
		}

		addr, err = wallet.GenerateNewAddress(name)
		if err != nil {
			errorLabel.SetText("could not generate new address, " + err.Error())
			errorLabel.Show()
			widget.Refresh(errorLabel)
			return
		}
		//if there was a rectified error and user clicks the generate again, this hides the error text
		if errorLabel.Hidden == false {
			errorLabel.Hide()
			widget.Refresh(errorLabel)
		}

		generatedAddress.SetText(addr)
		widget.Refresh(generatedAddress)

		png, _ := qrcode.Encode(addr, qrcode.High, 256)
		qrImage.Resource = fyne.NewStaticResource("Address", png)
		qrImage.Show()
		canvas.Refresh(qrImage)

		if generatedAddress.Hidden {
			generatedAddress.Show()
			copy.Show()
			widget.Refresh(generatedAddress)
			widget.Refresh(copy)
		}
	})
	button.Disable()

	output := widget.NewVBox(
		label,
		info,
		widget.NewHBox(accountLabel, receive.accountSelect),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(button.MinSize()), button),
		widgets.NewVSpacer(10),
		widget.NewHBox(layout.NewSpacer(), qrImage, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), generatedAddress, copy, layout.NewSpacer()),
		errorLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(10), output)
}

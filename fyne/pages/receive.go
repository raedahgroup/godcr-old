package pages

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/core"
	"github.com/raedahgroup/godcr/fyne/widgets"
	"github.com/skip2/go-qrcode"
)

func ReceivePageContent(wallet   *dcrlibwallet.LibWallet, window fyne.Window) fyne.CanvasObject {
	// if there were to be situations, wallet fails and new address cant be generated, then simply show fyne logo
	qrImage := canvas.NewImageFromResource(theme.FyneLogo())
	qrImage.SetMinSize(fyne.NewSize(300, 300))

	label := widget.NewLabelWithStyle("Receiving Funds", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	info := widget.NewLabelWithStyle("Each time you request a payment, a new address is created to protect your privacy.", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Italic: true})
	accountLabel := widget.NewLabelWithStyle("Account:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	var generatedAddress *widget.Label
	errorLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// todo: remove this after the hide bug on fyne is fixed
	// to test you can fmt.Println(errorLabel.Hidden) before the second hide function
	errorLabel.Hide()
	errorLabel.Hide()

	var addr string
	copyToClipboard := widget.NewToolbar(widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
		clipboard := window.Clipboard()
		clipboard.SetContent(addr)
	}))

	button := widget.NewButton("Generate Address", func() {
		name, err := wallet.AccountNumber(receive.accountSelect.Selected)
		if err != nil {
			errorLabel.SetText("error getting account name, " + err.Error())
			errorLabel.Show()
			return
		}

		addr, err = wallet.NextAddress(int32(name))
		if err != nil {
			errorLabel.SetText("could not generate new address, " + err.Error())
			errorLabel.Show()
			return
		}
		// if there was a rectified error and user clicks the generate again, this hides the error text
		if errorLabel.Hidden == false {
			errorLabel.Hide()
		}

		generatedAddress.SetText(addr)

		png, _ := qrcode.Encode(addr, qrcode.High, 256)
		qrImage.Resource = fyne.NewStaticResource("Address", png)
		qrImage.Show()
		canvas.Refresh(qrImage)
	})
	button.Disable()

	// get account and generate address on start
	accounts, err := core.AccountsOverview(wallet, core.DefaultRequiredConfirmations)
	if err != nil {
		errorLabel = widget.NewLabelWithStyle("Could not retrieve account information"+err.Error(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		// todo: log to file
		fmt.Println(err.Error())
		errorLabel.Show()
	}

	var options []string
	for _, account := range accounts {
		options = append(options, account.Name)
	}

	receive.accountSelect = widget.NewSelect(options, func(s string) {
		if button.Disabled() == true {
			button.Enable()
		}
	})

	receive.accountSelect.SetSelected(accounts[0].Name)

	addr, err = wallet.NextAddress(0)
	if err != nil {
		errorLabel.SetText("could not generate new address, " + err.Error())
		errorLabel.Show()
	}

	if errorLabel.Hidden {
		generatedAddress = widget.NewLabelWithStyle(addr, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		png, _ := qrcode.Encode(addr, qrcode.High, 256)
		qrImage.Resource = fyne.NewStaticResource("Address", png)
		canvas.Refresh(qrImage)
	}

	output := widget.NewVBox(
		label,
		info,
		widget.NewHBox(accountLabel, receive.accountSelect),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(button.MinSize()), button),
		widgets.NewVSpacer(10),
		widget.NewHBox(layout.NewSpacer(), qrImage, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), generatedAddress, copyToClipboard, layout.NewSpacer()),
		errorLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(10), output)
}

type receivePageData struct {
	accountSelect *widget.Select
}

var receive receivePageData

// todo: remove this when account page is implemented
func receivePageUpdates(wallet   *dcrlibwallet.LibWallet) {
	accounts, _ := core.AccountsOverview(wallet, core.DefaultRequiredConfirmations)

	var options []string
	for _, account := range accounts {
		options = append(options, account.Name)
	}
	receive.accountSelect.Options = options
	widget.Refresh(receive.accountSelect)
}

package pages

import (
	"fmt"

	"github.com/skip2/go-qrcode"

	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"

	"fyne.io/fyne/canvas"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

//todo: should we make concurrent checks if users add a new account?
func receivePage(wallet godcrApp.WalletMiddleware, window fyne.Window) fyne.CanvasObject {
	//initially load the fyne logo and hide it
	qrImage := canvas.NewImageFromResource(theme.FyneLogo())
	qrImage.SetMinSize(fyne.NewSize(200, 200))
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

	accountSelect := widget.NewSelect(options, func(s string) {
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
		name, err := wallet.AccountNumber(accountSelect.Selected)
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
		widget.NewHBox(accountLabel, accountSelect),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(button.MinSize()), button),
		widgets.NewVSpacer(10),
		widget.NewHBox(layout.NewSpacer(), qrImage, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), generatedAddress, copy, layout.NewSpacer()),
		errorLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(10), output)
}

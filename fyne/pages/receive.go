package pages

import (
	"fmt"
	"image/color"

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

	//errorLabel is defined here so as to be able to update its color when theme is changed
	errorLabel *canvas.Text
}

var receive receivePageData

func receivePageUpdates(wallet godcrApp.WalletMiddleware) {
	accounts, _ := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)

	var name []string
	for _, account := range accounts {
		name = append(name, account.Name)
	}

	name = append(name, "test")

	receive.accountSelect.Options = name
}

//todo: should we make concurrent checks if users add a new account?
func receivePage(wallet godcrApp.WalletMiddleware, window fyne.Window) fyne.CanvasObject {
	//if there were to be situations, wallet fails and new address cant be generated, then simply show fyne logo
	qrImage := canvas.NewImageFromResource(theme.FyneLogo())
	qrImage.SetMinSize(fyne.NewSize(300, 300))

	label := widget.NewLabelWithStyle("Receiving Funds", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	info := widget.NewLabelWithStyle("Each time you request a payment, a new address is created to protect your privacy.", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Italic: true})
	accountLabel := widget.NewLabelWithStyle("Account:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	generatedAddress := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	receive.errorLabel = canvas.NewText("", color.RGBA{255, 0, 0, 0})
	receive.errorLabel.Alignment = fyne.TextAlignCenter
	receive.errorLabel.TextStyle = fyne.TextStyle{Bold: true}
	receive.errorLabel.Hide()

	var addr string
	copy := widget.NewToolbar(widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
		clipboard := window.Clipboard()
		clipboard.SetContent(addr)
	}))

	button := widget.NewButton("Generate Address", func() {
		name, err := wallet.AccountNumber(receive.accountSelect.Selected)
		if err != nil {
			receive.errorLabel.Text = ("error getting account name, " + err.Error())
			receive.errorLabel.Show()
			canvas.Refresh(receive.errorLabel)
			return
		}

		addr, err = wallet.GenerateNewAddress(name)
		if err != nil {
			receive.errorLabel.Text = ("could not generate new address, " + err.Error())
			receive.errorLabel.Show()
			canvas.Refresh(receive.errorLabel)
			return
		}
		//if there was a rectified error and user clicks the generate again, this hides the error text
		if receive.errorLabel.Hidden == false {
			receive.errorLabel.Hide()
		}

		generatedAddress.SetText(addr)

		png, _ := qrcode.Encode(addr, qrcode.High, 256)
		qrImage.Resource = fyne.NewStaticResource("Address", png)
		qrImage.Show()
		canvas.Refresh(qrImage)
	})
	button.Disable()

	//get account and generate address on start
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		receive.errorLabel.Text = "Could not retrieve account information" + err.Error()
		//todo: log to file
		fmt.Println(err.Error())
		receive.errorLabel.Show()
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

	if receive.errorLabel.Hidden {
		addr, err = wallet.GenerateNewAddress(0)
		if err != nil {
			receive.errorLabel.Text = ("could not generate new address, " + err.Error())
			receive.errorLabel.Show()
		}

		if receive.errorLabel.Hidden {
			generatedAddress = widget.NewLabelWithStyle(addr, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

			png, _ := qrcode.Encode(addr, qrcode.High, 256)
			qrImage.Resource = fyne.NewStaticResource("Address", png)
			canvas.Refresh(qrImage)
		}
	}

	output := widget.NewVBox(
		label,
		info,
		widget.NewHBox(accountLabel, receive.accountSelect),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(button.MinSize()), button),
		widgets.NewVSpacer(10),
		widget.NewHBox(layout.NewSpacer(), qrImage, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), generatedAddress, copy, layout.NewSpacer()),
		receive.errorLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(10), output)
}

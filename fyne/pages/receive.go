package pages

import (
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

var receivePageContainer pageContainer

func receivePage(wallet godcrApp.WalletMiddleware, window fyne.Window) {
	// If there were to be situations, wallet fails and new address cant be generated, then simply show fyne logo.
	qrImage := canvas.NewImageFromResource(theme.FyneLogo())
	qrImage.SetMinSize(fyne.NewSize(300, 300))

	label := widget.NewLabelWithStyle("Receiving Funds", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	info := widget.NewLabelWithStyle("Each time you request a payment, a new address is created to protect your privacy.", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Italic: true})
	accountLabel := widget.NewLabelWithStyle("Account:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	generatedAddress := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	errorLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	errorLabel.Hide()
	errorLabel.Hide()

	var addr string
	copy := widget.NewToolbar(widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
		clipboard := window.Clipboard()
		clipboard.SetContent(addr)
	}))

	errorHandler := func(err string) {
		errorLabel.Show()
		errorLabel.SetText("error getting account name, " + err)
	}

	var accountSelect *widget.Select
	generateNewAddress := func() {
		name, err := wallet.AccountNumber(accountSelect.Selected)
		if err != nil {
			errorHandler(err.Error())
			return
		}

		addr, err = wallet.GenerateNewAddress(name)
		if err != nil {
			errorHandler(err.Error())
			return
		}
		widget.Refresh(generatedAddress)
		generatedAddress.SetText(addr)

		png, err := qrcode.Encode(addr, qrcode.High, 256)
		if err != nil {
			errorHandler(err.Error())
			return
		}
		// If there was a rectified error and user clicks the generate again, this hides the error text.
		if !errorLabel.Hidden {
			errorLabel.Hide()
		}

		qrImage.Resource = fyne.NewStaticResource("Address", png)
		qrImage.Show()
		canvas.Refresh(qrImage)
	}

	button := widget.NewButton("Generate Address", func() {
		generateNewAddress()
	})

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		receivePageContainer.container.Children = []fyne.CanvasObject{widget.NewLabel("error getting account name " + err.Error())}
		widget.Refresh(receivePageContainer.container)
		return
	}
	errorLabel.Hide()

	var options []string
	for _, account := range accounts {
		options = append(options, account.Name)
	}
	accountSelect = widget.NewSelect(options, func(selected string) {
		generateNewAddress()
	})
	accountSelect.SetSelected(accountSelect.Options[0])

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

	receivePageContainer.container.Children = widget.NewHBox(widgets.NewHSpacer(10), output).Children
	widget.Refresh(receivePageContainer.container)
}

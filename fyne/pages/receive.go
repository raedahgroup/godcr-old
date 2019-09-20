package pages

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/widgets"
	"github.com/skip2/go-qrcode"
)

func receivePageContent(dcrlw *dcrlibwallet.LibWallet, window fyne.Window) fyne.CanvasObject {
	icons, err := getIcons(collapse)
	if err != nil {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabel("Could not load image: "+err.Error()))
	}

	qrImage := canvas.NewImageFromResource(theme.FyneLogo())
	qrImage.SetMinSize(fyne.NewSize(300, 300))

	label := widget.NewLabelWithStyle("Receiving Funds", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	info := widget.NewLabelWithStyle("Each time you request a payment, a new address is created to protect your privacy.", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Italic: true})
	//accountLabel := widget.NewLabelWithStyle("Account:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
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

	var accountSelected string

	generateNewAddress := func() {
		name, err := dcrlw.AccountNumber(accountSelected)
		if err != nil {
			errorHandler(err.Error())
			return
		}

		addr, err = dcrlw.NextAddress(int32(name))
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

	// button := widget.NewButton("Generate Address", func() {
	// 	generateNewAddress()
	// })

	accounts, err := dcrlw.GetAccountsRaw(0) //wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return widget.NewLabel("error getting account name " + err.Error())
	}
	errorLabel.Hide()

	account := accounts.Acc[0]
	accountSelected = account.Name
	accountSelectedLabel := widget.NewLabel(accountSelected)
	accountSelectedBalanceLabel := widget.NewLabel(dcrutil.Amount(account.TotalBalance).String())
	generateNewAddress()

	var accountSelectionPopup *widget.PopUp

	accountsBox := widget.NewVBox()
	for i, account := range accounts.Acc {
		if account.Name == "imported" {
			continue
		}

		var accountName = account.Name
		var balance = dcrutil.Amount(account.Balance.Total).String()

		accountProperties := widget.NewHBox()
		accountProperties.Append(widgets.NewHSpacer(17))
		accountProperties.Append(widget.NewIcon(theme.ContentAddIcon()))
		accountProperties.Append(widgets.NewHSpacer(18))

		spendableLabel := canvas.NewText("Spendable", color.Black)
		spendableLabel.TextSize = 12
		spendableLabel.Alignment = fyne.TextAlignLeading

		spendableAmountLabel := canvas.NewText(dcrutil.Amount(account.Balance.Spendable).String(), color.Black)
		spendableAmountLabel.TextSize = 12
		spendableAmountLabel.Alignment = fyne.TextAlignCenter

		accountProperties.Append(widget.NewVBox(
			widget.NewLabel(accountName),
			spendableLabel,
		))
		accountProperties.Append(widgets.NewHSpacer(84))
		accountProperties.Append(widget.NewVBox(
			widget.NewLabel(balance),
			spendableAmountLabel,
		))

		accountProperties.Append(widgets.NewHSpacer(8))

		checkmarkIcon := widget.NewIcon(theme.ConfirmIcon())
		if i != 0 {
			checkmarkIcon.Hide()
		}
		accountProperties.Append(checkmarkIcon)

		accountsBox.Append(widgets.NewClickableBox(accountProperties, func() {
			// hide checkmark icon of other accounts
			for _, children := range accountsBox.Children {
				if box, ok := children.(*widgets.ClickableBox); !ok {
					continue
				} else {
					if len(box.Children) != 8 {
						fmt.Println("Not working")
						continue
					}
					fmt.Println("Works")
					if icon, ok := box.Children[7].(*widget.Icon); !ok {
						continue
					} else {
						icon.Hide()
						fmt.Println("Works")
					}
				}
			}

			checkmarkIcon.Show()
			accountSelectedLabel.SetText(accountName)
			accountSelectedBalanceLabel.SetText(balance)
			accountSelected = accountName
			generateNewAddress()
			accountSelectionPopup.Hide()
		}))
	}

	// create a popup that has account names with spendable amount
	accountSelectionPopup = widget.NewPopUp(accountsBox, window.Canvas())
	accountSelectionPopup.Hide()

	var accountPopup *widgets.ClickableIcon
	accountPopup = widgets.NewClickableIcon(icons[collapse], nil, func() {
		accountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(accountPopup).Add(fyne.NewPos(0, accountPopup.Size().Height)))
		accountSelectionPopup.Show()
	})

	accountTab := widget.NewHBox(widget.NewIcon(theme.ContentAddIcon()), widgets.NewHSpacer(16), accountSelectedLabel, widgets.NewHSpacer(76), accountSelectedBalanceLabel, widgets.NewHSpacer(8), accountPopup)

	output := widget.NewVBox(
		label,
		info,
		accountTab,
		widgets.NewVSpacer(10),
		widget.NewHBox(layout.NewSpacer(), qrImage, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), generatedAddress, copy, layout.NewSpacer()),
		errorLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(10), output)
}

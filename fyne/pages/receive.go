package pages

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
	"github.com/skip2/go-qrcode"
)

const receivingDecredHint = "Each time you request a\npayment, a new address is\ncreated to protect your privacy."

var receiveHandler struct {
	generatedReceiveAddress string
	recieveAddressError     error

	generatedQrCode []byte
	qrcodeError     error

	wallet              *dcrlibwallet.LibWallet
	accountNumber       uint32
	selectedAccountName string
}

func ReceivePageContent(dcrlw *dcrlibwallet.LibWallet, window fyne.Window, tabmenu *widget.TabContainer) fyne.CanvasObject {
	receiveHandler.wallet = dcrlw

	icons, err := assets.GetIcons(assets.CollapseIcon, assets.ReceiveAccountIcon, assets.MoreIcon, assets.InfoIcon)
	qrImage := widget.NewIcon(theme.FyneLogo())

	accountCopiedLabel := widget.NewLabelWithStyle("Address copied", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	widget.Refresh(accountCopiedLabel)
	accountCopiedLabel.Hide()

	label := widget.NewLabelWithStyle("Receiving Decred", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	generatedAddress := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	errorLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	errorLabel.Hide()

	accounts, err := receiveHandler.wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return widget.NewLabel(fmt.Sprintf("Error: %s", err.Error()))
	}

	errorHandler := func(err string) {
		errorLabel.SetText(err)
		errorLabel.Show()
	}

	generateAddress := func(generateNewAddress bool) {
		accountNumber, err := receiveHandler.wallet.AccountNumber(receiveHandler.selectedAccountName)
		if err != nil {
			errorHandler(fmt.Sprintf("Error: %s", err.Error()))
			return
		}
		receiveHandler.accountNumber = accountNumber

		if generateNewAddress {
			generateNewAddressAndQrCode()
		} else {
			generateAddressAndQrCode()
		}

		if receiveHandler.recieveAddressError != nil {
			errorHandler(fmt.Sprintf("Error: %s", receiveHandler.recieveAddressError.Error()))
			return
		}
		if receiveHandler.qrcodeError != nil {
			errorHandler(fmt.Sprintf("Error: %s", receiveHandler.qrcodeError.Error()))
			return
		}

		widget.Refresh(generatedAddress)
		generatedAddress.SetText(receiveHandler.generatedReceiveAddress)

		// If there was a rectified error and user clicks the generate address button again, this hides the error text.
		if !errorLabel.Hidden {
			errorLabel.Hide()
		}

		qrImage.SetResource(fyne.NewStaticResource("", receiveHandler.generatedQrCode))
	}

	// receiving decred hint-text pop-up
	var clickableInfoIcon *widgets.ClickableIcon
	clickableInfoIcon = widgets.NewClickableIcon(icons[assets.InfoIcon], nil, func() {
		label := widget.NewLabelWithStyle(receivingDecredHint, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
		gotItLabel := canvas.NewText("Got it", color.RGBA{41, 112, 255, 255})
		gotItLabel.TextStyle = fyne.TextStyle{Bold: true}
		gotItLabel.TextSize = 16

		var popup *widget.PopUp
		popup = widget.NewPopUp(widget.NewVBox(
			widgets.NewVSpacer(24),
			widget.NewHBox(widgets.NewHSpacer(24), widget.NewLabelWithStyle("Receiving Decred", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
			widgets.NewVSpacer(10),
			widget.NewHBox(widgets.NewHSpacer(24), label),
			widgets.NewVSpacer(10),
			widget.NewHBox(layout.NewSpacer(), widgets.NewClickableBox(widget.NewVBox(gotItLabel), func() { popup.Hide() }), widgets.NewHSpacer(24)),
			widgets.NewVSpacer(18),
		), window.Canvas())

		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(clickableInfoIcon).Add(fyne.NewPos(0, clickableInfoIcon.Size().Height)))
	})

	// generate new address pop-up
	var clickableMoreIcon *widgets.ClickableIcon
	clickableMoreIcon = widgets.NewClickableIcon(icons[assets.MoreIcon], nil, func() {
		var popup *widget.PopUp
		popup = widget.NewPopUp(widgets.NewClickableBox(widget.NewHBox(widget.NewLabel("Generate new address")), func() {
			generateAddress(true)
			popup.Hide()
		}), window.Canvas())

		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(clickableMoreIcon).Add(fyne.NewPos(0, clickableMoreIcon.Size().Height)))
		popup.Show()
	})

	// automatically generate address for first account
	account := accounts.Acc[0]
	receiveHandler.selectedAccountName = account.Name
	generateAddress(false)

	selectedAccountLabel := widget.NewLabel(receiveHandler.selectedAccountName)
	selectedAccountBalanceLabel := widget.NewLabel(dcrutil.Amount(account.TotalBalance).String())

	var accountSelectionPopup *widget.PopUp
	accountsBox := widget.NewVBox()

	for i, account := range accounts.Acc {
		if account.Name == "imported" {
			continue
		}

		var accountName = account.Name
		var balance = dcrutil.Amount(account.Balance.Total).String()

		accountProperties := widget.NewHBox()
		accountProperties.Append(widgets.NewHSpacer(15))
		accountProperties.Append(widget.NewIcon(icons[assets.ReceiveAccountIcon]))
		accountProperties.Append(widgets.NewHSpacer(15))

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
						continue
					}

					if icon, ok := box.Children[7].(*widget.Icon); !ok {
						continue
					} else {
						icon.Hide()
					}
				}
			}

			// checkmarkIcon.Show()
			selectedAccountLabel.SetText(accountName)
			selectedAccountBalanceLabel.SetText(balance)
			receiveHandler.selectedAccountName = accountName
			generateAddress(false)
			accountSelectionPopup.Hide()
		}))
	}

	hidePopUp := widget.NewHBox(widgets.NewHSpacer(16),
		widgets.NewClickableIcon(theme.CancelIcon(), nil, func() { accountSelectionPopup.Hide() }),
		widgets.NewHSpacer(16), widget.NewLabelWithStyle("Receiving account", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer())

	// accountSelectionPopup create a popup that has account names with spendable amount
	accountSelectionPopup = widget.NewPopUp(widget.NewVBox(
		widgets.NewVSpacer(5), hidePopUp, widgets.NewVSpacer(5), canvas.NewLine(color.Black), accountsBox),
		window.Canvas())
	accountSelectionPopup.Hide()

	// accountTab shows the selected account
	accountTab := widget.NewHBox(widget.NewIcon(icons[assets.ReceiveAccountIcon]), widgets.NewHSpacer(16),
		selectedAccountLabel, widgets.NewHSpacer(76), selectedAccountBalanceLabel, widgets.NewHSpacer(8), widget.NewIcon(icons[assets.CollapseIcon]))

	var accountDropdown *widgets.ClickableBox
	accountDropdown = widgets.NewClickableBox(accountTab, func() {
		accountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			accountDropdown).Add(fyne.NewPos(0, accountDropdown.Size().Height)))
		accountSelectionPopup.Show()
	})

	// copyAddressAction enables address copying
	copyAddressAction := widgets.NewClickableBox(widget.NewVBox(widget.NewLabelWithStyle("(Tap to copy)", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})),
		func() {
			clipboard := window.Clipboard()
			clipboard.SetContent(receiveHandler.generatedReceiveAddress)

			accountCopiedLabel.Show()
			// only hide accountCopiedLabel text if user is currently on the page after 5secs
			if accountCopiedLabel.Hidden == false {
				time.AfterFunc(time.Second*10, func() {
					if tabmenu.CurrentTabIndex() == 3 {
						accountCopiedLabel.Hide()
					}
				})
			}
		})

	output := widget.NewVBox(
		widgets.NewVSpacer(10),
		accountCopiedLabel,
		widget.NewHBox(label, widgets.NewHSpacer(110), clickableInfoIcon, widgets.NewHSpacer(26), clickableMoreIcon),
		widgets.NewVSpacer(18),
		accountDropdown,
		widgets.NewVSpacer(32),
		// due to original width content on page being small layout spacing isn't efficient
		// therefore requiring an additional spacing by a 100 width
		widget.NewHBox(layout.NewSpacer(), widgets.NewHSpacer(100), fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(300, 300)), qrImage), layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), widgets.NewHSpacer(100), generatedAddress, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), widgets.NewHSpacer(100), copyAddressAction, layout.NewSpacer()),
		errorLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(18), output)
}

func generateAddressAndQrCode() {
	receiveHandler.generatedReceiveAddress, receiveHandler.recieveAddressError = receiveHandler.wallet.CurrentAddress(int32(receiveHandler.accountNumber))
	receiveHandler.generatedQrCode, receiveHandler.qrcodeError = generateQrcode(receiveHandler.generatedReceiveAddress)
}

func generateNewAddressAndQrCode() {
	receiveHandler.generatedReceiveAddress, receiveHandler.recieveAddressError = receiveHandler.wallet.NextAddress(int32(receiveHandler.accountNumber))
	receiveHandler.generatedQrCode, receiveHandler.qrcodeError = generateQrcode(receiveHandler.generatedReceiveAddress)
}

func generateQrcode(generatedReceiveAddress string) ([]byte, error) {
	// generate qrcode
	png, err := qrcode.Encode(generatedReceiveAddress, qrcode.High, 256)
	if err != nil {
		return nil, err
	}

	return png, nil
}

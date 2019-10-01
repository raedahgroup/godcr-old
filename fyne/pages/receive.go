package pages

import (
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/fyne/widgets"
	"github.com/skip2/go-qrcode"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/utils"
)

const (
	defaultRequiredConfirmations = 2
)
func ReceivePageContent(dcrlw *dcrlibwallet.LibWallet, window fyne.Window, tabmenu *widget.TabContainer) fyne.CanvasObject {
	icons, err := utils.GetIcons(utils.Collapse, utils.ReceiveAccount, utils.More, utils.Info)
	qrImage := widget.NewIcon(theme.FyneLogo())

	accountCopiedLabel := widget.NewLabelWithStyle("Address copied", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	widget.Refresh(accountCopiedLabel)
	accountCopiedLabel.Hide()

	label := widget.NewLabelWithStyle("Receive", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	generatedAddress := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	errorLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	errorLabel.Hide()

	var addr string

	errorHandler := func(err string) {
		errorLabel.SetText("error getting account name, " + err)
		errorLabel.Show()
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

		// If there was a rectified error and user clicks the generate address button again, this hides the error text.
		if !errorLabel.Hidden {
			errorLabel.Hide()
		}

		qrImage.SetResource(fyne.NewStaticResource("", png))
	}

	var clickableInfoIcon *widgets.ClickableIcon

	clickableInfoIcon = widgets.NewClickableIcon(icons[utils.Info], nil, func() {
		label := widget.NewLabelWithStyle("Each time you request a\npayment, a new address is\ncreated to protect your privacy.", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
		gotItLabel := canvas.NewText("Got it", color.RGBA{41, 112, 255, 255})
		gotItLabel.TextStyle = fyne.TextStyle{Bold: true}
		gotItLabel.TextSize = 16

		var popup *widget.PopUp
		popup = widget.NewPopUp(widget.NewVBox(
			widgets.NewVSpacer(24),
			widget.NewHBox(widgets.NewHSpacer(24), widget.NewLabelWithStyle("Receive DCR", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
			widgets.NewVSpacer(10),
			widget.NewHBox(widgets.NewHSpacer(24), label),
			widgets.NewVSpacer(10),
			widget.NewHBox(layout.NewSpacer(), widgets.NewClickableBox(widget.NewVBox(gotItLabel), func() { popup.Hide() }), widgets.NewHSpacer(24)),
			widgets.NewVSpacer(18),
		), window.Canvas())

		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(clickableInfoIcon).Add(fyne.NewPos(0, clickableInfoIcon.Size().Height)))
	})

	var clickableMoreIcon *widgets.ClickableIcon

	clickableMoreIcon = widgets.NewClickableIcon(icons[utils.More], nil, func() {
		var popup *widget.PopUp
		popup = widget.NewPopUp(widgets.NewClickableBox(widget.NewHBox(widget.NewLabel("Generate new address")), func() {
			generateNewAddress()
			popup.Hide()
		}), window.Canvas())

		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(clickableMoreIcon).Add(fyne.NewPos(0, clickableMoreIcon.Size().Height)))
		popup.Show()
	})

	accounts, err := dcrlw.GetAccountsRaw(defaultRequiredConfirmations)
	if err != nil {
		return widget.NewLabel("error getting account name " + err.Error())
	}

	// automatically generate address for first account
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
		accountProperties.Append(widget.NewIcon(icons[utils.ReceiveAccount]))
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
						continue
					}

					if icon, ok := box.Children[7].(*widget.Icon); !ok {
						continue
					} else {
						icon.Hide()
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

	hidePopUp := widget.NewHBox(widgets.NewHSpacer(16),
		widgets.NewClickableIcon(theme.CancelIcon(), nil, func() { accountSelectionPopup.Hide() }),
		widgets.NewHSpacer(16), widget.NewLabelWithStyle("Receiving account", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer())

	// create a popup that has account names with spendable amount
	accountSelectionPopup = widget.NewPopUp(widget.NewVBox(
		widgets.NewVSpacer(5), hidePopUp, widgets.NewVSpacer(5), canvas.NewLine(color.Black), accountsBox),
		window.Canvas())
	accountSelectionPopup.Hide()

	// accountTab shows the selected account
	accountTab := widget.NewHBox(widget.NewIcon(icons[utils.ReceiveAccount]), widgets.NewHSpacer(16),
		accountSelectedLabel, widgets.NewHSpacer(76), accountSelectedBalanceLabel, widgets.NewHSpacer(8), widget.NewIcon(icons[utils.Collapse]))

	var accountDropdown *widgets.ClickableBox
	accountDropdown = widgets.NewClickableBox(accountTab, func() {
		accountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			accountDropdown).Add(fyne.NewPos(0, accountDropdown.Size().Height)))
		accountSelectionPopup.Show()
	})

	copyAddressAction := widgets.NewClickableBox(widget.NewVBox(widget.NewLabelWithStyle("(Tap to copy)", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})),
		func() {
			clipboard := window.Clipboard()
			clipboard.SetContent(addr)

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

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

var receivePageDataHandler struct {
	selectedAccountName     string
	generatedReceiveAddress string
	receiveAddressError     error
	wallet                  *dcrlibwallet.LibWallet
	generatedAddressLabel   *widget.Label
	errorLabel              *widget.Label
	qrImage                 *widget.Icon
}

func ReceivePageContent(dcrlw *dcrlibwallet.LibWallet, window fyne.Window, tabmenu *widget.TabContainer) fyne.CanvasObject {
	receivePageDataHandler.wallet = dcrlw

	// error handler
	receivePageDataHandler.errorLabel = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	receivePageDataHandler.errorLabel.Hide()

	receivePageDataHandler.qrImage = widget.NewIcon(theme.FyneLogo())
	receivePageDataHandler.generatedAddressLabel = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	const pageTitleText = "Receive DCR"
	const receivingDecredHint = "Each time you request a payment, a new \naddress is created to protect your privacy."

	// get icons used on the page.
	icons, err := assets.GetIcons(assets.CollapseIcon, assets.ReceiveAccountIcon, assets.MoreIcon, assets.InfoIcon)

	// clickableInfoIcon holds receiving decred hint-text pop-up
	var clickableInfoIcon *widgets.ClickableIcon
	clickableInfoIcon = widgets.NewClickableIcon(icons[assets.InfoIcon], nil, func() {
		receivingDecredHintLabel := widget.NewLabelWithStyle(receivingDecredHint, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
		gotItLabel := canvas.NewText("Got it", color.RGBA{41, 112, 255, 255})
		gotItLabel.TextStyle = fyne.TextStyle{Bold: true}
		gotItLabel.TextSize = 16

		var receivingDecredHintLabelPopUp *widget.PopUp
		receivingDecredHintLabelPopUp = widget.NewPopUp(widget.NewVBox(
			widgets.NewVSpacer(14),
			widget.NewHBox(widgets.NewHSpacer(14), widget.NewLabelWithStyle(pageTitleText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
			widgets.NewVSpacer(5),
			widget.NewHBox(widgets.NewHSpacer(14), receivingDecredHintLabel, widgets.NewHSpacer(14)),
			widgets.NewVSpacer(10),
			widget.NewHBox(layout.NewSpacer(), widgets.NewClickableBox(widget.NewVBox(gotItLabel), func() { receivingDecredHintLabelPopUp.Hide() }), widgets.NewHSpacer(14)),
			widgets.NewVSpacer(14),
		), window.Canvas())

		receivingDecredHintLabelPopUp.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(clickableInfoIcon).Add(fyne.NewPos(0, clickableInfoIcon.Size().Height)))
	})

	// generate new address pop-up
	var clickableMoreIcon *widgets.ClickableIcon
	clickableMoreIcon = widgets.NewClickableIcon(icons[assets.MoreIcon], nil, func() {
		var generateNewAddressPopup *widget.PopUp
		generateNewAddressPopup = widget.NewPopUp(widgets.NewClickableBox(widget.NewHBox(widget.NewLabel("Generate new address")), func() {
			generateAddressAndQrCode(true)
			generateNewAddressPopup.Hide()
		}), window.Canvas())

		generateNewAddressPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(clickableMoreIcon).Add(fyne.NewPos(0, clickableMoreIcon.Size().Height)))
		generateNewAddressPopup.Show()
	})

	// get user accounts
	accounts, err := receivePageDataHandler.wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return widget.NewLabel(fmt.Sprintf("Error: %s", err.Error()))
	}

	// automatically generate address for first account
	receivePageDataHandler.selectedAccountName = accounts.Acc[0].Name
	selectedAccountLabel := widget.NewLabel(receivePageDataHandler.selectedAccountName)
	selectedAccountBalanceLabel := widget.NewLabel(dcrutil.Amount(accounts.Acc[0].TotalBalance).String())
	generateAddressAndQrCode(false)

	var accountSelectionPopup *widget.PopUp
	accountListWidget := widget.NewVBox()

	for _, account := range accounts.Acc {
		if account.Name == "imported" {
			continue
		}

		spendableLabel := canvas.NewText("Spendable", color.White)
		spendableLabel.TextSize = 10
		spendableLabel.Alignment = fyne.TextAlignCenter

		spendableAmountLabel := canvas.NewText(dcrutil.Amount(account.Balance.Spendable).String(), color.White)
		spendableAmountLabel.TextSize = 10
		spendableAmountLabel.Alignment = fyne.TextAlignCenter

		accountName := account.Name
		accountNameBox := widget.NewVBox(
			widget.NewLabel(accountName),
			spendableLabel,
		)

		accountBalance := dcrutil.Amount(account.Balance.Total).String()
		accountBalanceBox := widget.NewVBox(
			widget.NewLabel(accountBalance),
			spendableAmountLabel,
		)

		accountsView := widget.NewHBox(
			widgets.NewHSpacer(15),
			widget.NewIcon(icons[assets.ReceiveAccountIcon]),
			widgets.NewHSpacer(20),
			accountNameBox,
			widgets.NewHSpacer(20),
			accountBalanceBox,
			widgets.NewHSpacer(15),
		)

		accountListWidget.Append(widgets.NewClickableBox(accountsView, func() {
			selectedAccountLabel.SetText(accountName)
			selectedAccountBalanceLabel.SetText(accountBalance)
			receivePageDataHandler.selectedAccountName = accountName
			generateAddressAndQrCode(false)
			accountSelectionPopup.Hide()
		}))
	}

	accountSelectionPopupHeader := widget.NewHBox(
		widgets.NewHSpacer(16),
		widgets.NewClickableIcon(theme.CancelIcon(), nil, func() { accountSelectionPopup.Hide() }),
		widgets.NewHSpacer(16),
		widget.NewLabelWithStyle("Receiving account", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
	)

	// accountSelectionPopup create a popup that has account names with spendable amount
	accountSelectionPopup = widget.NewPopUp(
		widget.NewVBox(
			widgets.NewVSpacer(5),
			accountSelectionPopupHeader,
			widgets.NewVSpacer(5),
			canvas.NewLine(color.Black),
			accountListWidget,
			widgets.NewVSpacer(5),
		), window.Canvas(),
	)
	accountSelectionPopup.Hide()

	// accountTab shows the selected account
	accountTab := widget.NewHBox(
		widget.NewIcon(icons[assets.ReceiveAccountIcon]),
		widgets.NewHSpacer(15),
		selectedAccountLabel,
		widgets.NewHSpacer(30),
		selectedAccountBalanceLabel,
		widgets.NewHSpacer(8),
		widget.NewIcon(icons[assets.CollapseIcon]),
	)

	var accountDropdown *widgets.ClickableBox
	accountDropdown = widgets.NewClickableBox(accountTab, func() {
		accountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			accountDropdown).Add(fyne.NewPos(0, accountDropdown.Size().Height)))
		accountSelectionPopup.Show()
	})

	accountCopiedLabel := widget.NewLabelWithStyle("Address copied", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	accountCopiedLabel.Hide()

	// copyAddressAction enables address copying
	copyAddressAction := widgets.NewClickableBox(widget.NewVBox(widget.NewLabelWithStyle("(Tap to copy)", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})),
		func() {
			clipboard := window.Clipboard()
			clipboard.SetContent(receivePageDataHandler.generatedReceiveAddress)

			accountCopiedLabel.Show()
			widget.Refresh(accountCopiedLabel)

			// only hide accountCopiedLabel text if user is currently on the page after 2secs
			if accountCopiedLabel.Hidden == false {
				time.AfterFunc(time.Second*2, func() {
					if tabmenu.CurrentTabIndex() == 3 {
						accountCopiedLabel.Hide()
					}
				})
			}
		},
	)

	pageTitleLabel := widget.NewLabelWithStyle(pageTitleText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	output := widget.NewVBox(
		widgets.NewVSpacer(5),
		widget.NewHBox(pageTitleLabel, widgets.NewHSpacer(110), clickableInfoIcon, widgets.NewHSpacer(26), clickableMoreIcon),
		accountDropdown,
		widget.NewHBox(layout.NewSpacer(), widgets.NewHSpacer(70), accountCopiedLabel, layout.NewSpacer()),
		// due to original width content on page being small layout spacing isn't efficient
		// therefore requiring an additional spacing by a 100 width
		widget.NewHBox(layout.NewSpacer(), widgets.NewHSpacer(70), fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(300, 300)), receivePageDataHandler.qrImage), layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), widgets.NewHSpacer(70), receivePageDataHandler.generatedAddressLabel, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), widgets.NewHSpacer(70), copyAddressAction, layout.NewSpacer()),
		receivePageDataHandler.errorLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(18), output)
}

func generateAddressAndQrCode(generateNewAddress bool) {
	accountNumber, err := receivePageDataHandler.wallet.AccountNumber(receivePageDataHandler.selectedAccountName)
	if err != nil {
		errorHandler(fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	if generateNewAddress {
		receivePageDataHandler.generatedReceiveAddress, receivePageDataHandler.receiveAddressError = receivePageDataHandler.wallet.NextAddress(int32(accountNumber))
	} else {
		receivePageDataHandler.generatedReceiveAddress, receivePageDataHandler.receiveAddressError = receivePageDataHandler.wallet.CurrentAddress(int32(accountNumber))
	}

	if receivePageDataHandler.receiveAddressError != nil {
		errorHandler(fmt.Sprintf("Error: %s", receivePageDataHandler.receiveAddressError.Error()))
		return
	}

	widget.Refresh(receivePageDataHandler.generatedAddressLabel)
	generateQrCode(receivePageDataHandler.generatedReceiveAddress)
	receivePageDataHandler.generatedAddressLabel.SetText(receivePageDataHandler.generatedReceiveAddress)

	// If there was a rectified error and user clicks the generate address button again, this hides the error text.
	if !receivePageDataHandler.errorLabel.Hidden {
		receivePageDataHandler.errorLabel.Hide()
	}
}

func generateQrCode(generatedAddress string) {
	generatedQrCode, err := qrcode.Encode(generatedAddress, qrcode.High, 256)
	if err != nil {
		errorHandler(fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	receivePageDataHandler.qrImage.SetResource(fyne.NewStaticResource("", generatedQrCode))
}

func errorHandler(err string) {
	receivePageDataHandler.errorLabel.SetText(err)
	receivePageDataHandler.errorLabel.Show()
}

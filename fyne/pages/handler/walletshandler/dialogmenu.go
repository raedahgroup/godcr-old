package walletshandler

import (
	"encoding/hex"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (walletPage *WalletPageObject) dialogMenu(walletLabel *canvas.Text, posOfIcon fyne.Position, walletID int) *widget.PopUp {
	var popUp *widget.PopUp

	clickableText := func(text string, callFunc func()) *widgets.ClickableWidget {
		TextWithPadding := widget.NewHBox(widgets.NewHSpacer(values.SpacerSize12), widgets.NewTextWithSize(text, values.DefaultTextColor, 14), layout.NewSpacer(), widgets.NewHSpacer(values.SpacerSize40))
		textBox := widget.NewVBox(
			widgets.NewVSpacer(values.SpacerSize12),
			TextWithPadding,
			widgets.NewVSpacer(values.SpacerSize12),
		)

		return widgets.NewClickableWidget(textBox, callFunc)
	}
	wallet := walletPage.MultiWallet.WalletWithID(walletID)

	callFunc := func() {
		popUp.Hide()
	}

	renameWalletFunc := func() {
		walletPage.renameWalletPopUp(walletID, walletLabel)
	}

	dialogBox := widget.NewVBox(
		widgets.NewHSpacer(values.SpacerSize4),
		clickableText(values.SignMessage, func() { walletPage.signMessagePopUp(wallet, popUp) }),
		clickableText(values.VerifyMessage, callFunc),
		widgets.NewHSpacer(values.SpacerSize4),
		canvas.NewLine(values.StrippedLineColor),
		widgets.NewHSpacer(values.SpacerSize4),
		clickableText(values.ViewProperty, callFunc),
		widgets.NewHSpacer(values.SpacerSize4),
		canvas.NewLine(values.StrippedLineColor),
		widgets.NewHSpacer(values.SpacerSize4),
		clickableText(values.RenameWallet, renameWalletFunc),
		clickableText(values.WalletSettings, callFunc),
		widgets.NewHSpacer(values.SpacerSize4),
	)

	posX := dialogBox.MinSize().Width

	popUp = widget.NewPopUpAtPosition(dialogBox, walletPage.Window.Canvas(), posOfIcon.Subtract(fyne.NewPos(posX, 0).Subtract(fyne.NewPos(0, 20))))
	return popUp
}

func (walletPage *WalletPageObject) renameWalletPopUp(walletID int, walletLabel *canvas.Text) { //baseText string, onRename func(string) error, onCancel func(*widget.PopUp), otherCallFunc func(string)) {
	onRename := func(value string) error {
		return walletPage.MultiWallet.RenameWallet(walletID, value)
	}
	onCancel := func(popup *widget.PopUp) {
		popup.Hide()
	}
	otherCallFunc := func(value string) {
		walletLabel.Text = value
		walletPage.showLabel("Wallet renamed", walletPage.successLabel)
	}

	walletPage.renameAccountOrWalletPopUp(values.RenameWallet, values.RenameWalletPlaceHolder, onRename, onCancel, otherCallFunc)
}

func (walletPage *WalletPageObject) signMessagePopUp(wallet *dcrlibwallet.Wallet, dialogPopup *widget.PopUp) {
	dialogPopup.Hide()
	var stringedMessage string
	var maxResize fyne.Size
	var scrollableMessageBox *fyne.Container

	var popup *widget.PopUp
	successLabel := widgets.NewBorderedText("", fyne.NewSize(20, 16), values.Green)
	successLabel.Container.Hide()
	errorLabel := widgets.NewTextWithSize("", values.ErrorColor, 12)

	backIcon := widgets.NewImageButton(theme.NavigateBackIcon(), nil, func() {
		popup.Hide()
	})

	infoIcon := widgets.NewImageButton(walletPage.icons[assets.InfoIcon], nil, func() {
		var infoPopUp *widget.PopUp

		gotItText := canvas.NewText("Got it", values.Blue)
		gotItText.TextStyle.Bold = true

		gotItButton := widgets.NewClickableWidget(widget.NewHBox(gotItText), func() {
			infoPopUp.Hide()
			popup.Show()
		})

		infoDetails := widget.NewVBox(
			widgets.NewVSpacer(values.SpacerSize20),
			widgets.NewTextWithStyle(values.SignMessage, values.DefaultTextColor, fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, 18),
			widgets.NewVSpacer(values.SpacerSize12),
			canvas.NewText("Signing message with an", values.SignMessageBaseLabelColor),
			canvas.NewText("address' private key allows you to", values.SignMessageBaseLabelColor),
			canvas.NewText("prove that you are the owner of a", values.SignMessageBaseLabelColor),
			canvas.NewText("given address to a possible", values.SignMessageBaseLabelColor),
			canvas.NewText("counterparty.", values.SignMessageBaseLabelColor),
			widget.NewHBox(layout.NewSpacer(), gotItButton),
			widgets.NewVSpacer(values.SpacerSize20),
		)

		infoPopUp = widget.NewModalPopUp(widget.NewHBox(widgets.NewHSpacer(values.SpacerSize20), infoDetails, widgets.NewHSpacer(values.SpacerSize20)), walletPage.Window.Canvas())
	})

	label := widgets.NewTextWithSize(values.SignMessage, values.DefaultTextColor, 20)
	baseLabel := canvas.NewText(values.SignMessageBaseLabel, values.SignMessageBaseLabelColor)

	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder(values.AddressPlaceHolder)
	addressErrorLabel := widgets.NewTextWithSize("", values.ErrorColor, 12)

	messageEntry := widget.NewMultiLineEntry()
	messageEntry.SetPlaceHolder(values.MessagePlaceHolder)

	clearAllText := canvas.NewText(values.ClearAll, values.DisabledButtonColor)
	clearAllText.TextStyle.Bold = true
	clearAllButton := widgets.NewClickableWidget(widget.NewHBox(clearAllText), func() {
		addressEntry.SetText("")
		messageEntry.SetText("")
	})
	clearAllButton.Disable()

	var signButton *widgets.Button

	messageEntry.OnChanged = func(value string) {
		if value == "" && addressEntry.Text == "" {
			clearAllText.Color = values.DisabledButtonColor
			clearAllButton.Disable()
			clearAllText.Refresh()
			walletPage.WalletPageContents.Refresh()
			return
		}

		if addressErrorLabel.Hidden && addressEntry.Text != "" && signButton.Disabled() {
			signButton.Enable()
		}

		clearAllText.Color = values.Blue
		clearAllText.Refresh()
		clearAllButton.Enable()

		walletPage.WalletPageContents.Refresh()
	}

	addressEntry.OnChanged = func(value string) {
		if value == "" && messageEntry.Text == "" {
			clearAllText.Color = values.DisabledButtonColor
			clearAllButton.Disable()
			clearAllText.Refresh()
			signButton.Disable()
			signButton.Container.Refresh()
			addressErrorLabel.Hide()

			return
		}

		clearAllText.Color = values.Blue
		clearAllText.Refresh()
		clearAllButton.Enable()

		if value == "" && !addressErrorLabel.Hidden {
			addressErrorLabel.Hide()
			walletPage.WalletPageContents.Refresh()
			return
		}

		if wallet.IsAddressValid(value) {
			if wallet.HaveAddress(value) {
				addressErrorLabel.Hide()
				addressErrorLabel.Refresh()
				signButton.Enable()
				signButton.Container.Refresh()
				walletPage.WalletPageContents.Refresh()
				return
			}

			addressErrorLabel.Text = "Address does not belong to wallet"
			addressErrorLabel.Show()
			addressErrorLabel.Refresh()
			signButton.Disable()

		} else {
			addressErrorLabel.Text = "Not a valid address."
			addressErrorLabel.Show()
			addressErrorLabel.Refresh()
			signButton.Disable()
		}

		walletPage.WalletPageContents.Refresh()
	}

	signatureEntry := widget.NewEntry()
	signatureEntry.Disable()

	copyButton := widgets.NewButton(values.Blue, values.Copy, func() {
		walletPage.Window.Clipboard().SetContent(stringedMessage)
		walletPage.showLabel("Signature copied", successLabel)
	})
	copyButton.SetTextStyle(fyne.TextStyle{Bold: true})
	copyButton.SetMinSize(copyButton.MinSize().Add(fyne.NewSize(48, 24)))

	signatureEntryBox := widget.NewVBox(
		canvas.NewLine(values.StrippedLineColor),
		widgets.NewVSpacer(values.SpacerSize12),
		signatureEntry,
		widgets.NewVSpacer(values.SpacerSize12),
		widget.NewHBox(layout.NewSpacer(), copyButton.Container),
		widgets.NewVSpacer(values.SpacerSize12),
	)

	signButton = widgets.NewButton(values.Blue, values.Sign, func() {
		onConfirm := func(password string) error {
			message, err := wallet.SignMessage([]byte(password), addressEntry.Text, messageEntry.Text)
			if err != nil {
				return err
			}

			stringedMessage = hex.EncodeToString(message)
			var splittedWords string
			for i := 0; i < len(stringedMessage); i += 40 {
				if len(stringedMessage) > i+40 {
					splittedWords += stringedMessage[i : i+40]
					splittedWords += "\n"
				} else {
					splittedWords += stringedMessage[i:]
				}
			}
			signatureEntry.SetText(splittedWords)
			signButton.Disable()
			return nil
		}
		onError := func(err error) {
			errorLabel.Text = err.Error()
			errorLabel.Show()
			errorLabel.Refresh()
			walletPage.WalletPageContents.Refresh()

			popup.Show()
		}
		extraCalls := func() {
			popup.Show()
			walletPage.showLabel("Message signed", successLabel)
			signatureEntryBox.Show()

			scrollableMessageBox.Layout = layout.NewFixedGridLayout(maxResize)
			scrollableMessageBox.Refresh()
		}
		onCancel := func() {
			popup.Show()
		}

		passwordPopUp := multipagecomponents.PasswordPopUpObjects{
			onConfirm, onError, onCancel, extraCalls, values.ConfirmToSign, walletPage.Window,
		}
		passwordPopUp.PasswordPopUp()

	})
	signButton.SetTextStyle(fyne.TextStyle{Bold: true})
	signButton.SetMinSize(signButton.MinSize().Add(fyne.NewSize(48, 24)))
	signButton.Disable()

	signMessageBox := widget.NewHBox(widgets.NewHSpacer(values.SpacerSize20),
		widget.NewVBox(
			widgets.NewVSpacer(values.SpacerSize14),
			widget.NewHBox(backIcon, widgets.NewHSpacer(values.SpacerSize12), label, layout.NewSpacer(), infoIcon),
			widgets.NewVSpacer(values.SpacerSize4),
			successLabel.Container,
			widgets.NewVSpacer(values.SpacerSize4),
			baseLabel,
			widgets.NewVSpacer(values.SpacerSize4),
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(widget.NewLabel(values.TestAddress).MinSize().Add(fyne.NewSize(0, 10))), addressEntry),
			addressErrorLabel,
			widgets.NewVSpacer(values.SpacerSize12),
			messageEntry,
			widgets.NewVSpacer(values.SpacerSize12),
			widget.NewHBox(layout.NewSpacer(), widgets.CenterObject(clearAllButton, false), widgets.NewHSpacer(values.SpacerSize20), signButton.Container),

			widgets.NewVSpacer(values.SpacerSize12),
			signatureEntryBox,
		),
		widgets.NewHSpacer(values.SpacerSize20))

	maxResize = signMessageBox.MinSize()
	signatureEntryBox.Hide()
	scrollableMessageBox = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(signMessageBox.MinSize()), widget.NewScrollContainer(signMessageBox))

	popup = widget.NewModalPopUp(scrollableMessageBox, dialogPopup.Canvas)
}

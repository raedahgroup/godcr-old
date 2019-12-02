package sendpagehandler

import (
	"fmt"
	"image/color"
	"log"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/layouts"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func ConfirmationWindow(amountEntry, destinationAddressEntry *widget.Entry, downArrow, alert, reveal, conceal *fyne.StaticResource,
	window fyne.Window, selectedWalletName, sendingToSelfSelectedWalletName string, totalCostText, transactionFeeText, balanceAfterSendText,
	sendingSelectedAccountText, selfSendingSelectedAccountText string, sendingToSelf bool, transactionAuthor *dcrlibwallet.TxAuthor,
	showSuccess *widgets.Button, contents *widget.Box) {

	var confirmationPagePopup *widget.PopUp

	confirmLabel := canvas.NewText("Confirm to send", color.Black)
	confirmLabel.TextStyle.Bold = true
	confirmLabel.TextSize = 20

	errorLabel := canvas.NewText("Failed to send. Please try again.", color.White)
	errorLabel.Alignment = fyne.TextAlignCenter
	errorBar := canvas.NewRectangle(color.RGBA{237, 109, 71, 255})
	errorBar.SetMinSize(errorLabel.MinSize().Add(fyne.NewSize(20, 16)))

	errorLabelContainer := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), errorBar, errorLabel)
	errorLabelContainer.Hide()

	accountSelectionPopupHeader := widget.NewHBox(
		widgets.NewImageButton(theme.CancelIcon(), nil, func() { confirmationPagePopup.Hide() }),
		widgets.NewHSpacer(9),
		confirmLabel,
		widgets.NewHSpacer(170),
	)
	sendingSelectedWalletLabel := widget.NewLabelWithStyle(fmt.Sprintf("%s (%s)",
		sendingSelectedAccountText, selectedWalletName), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	trailingDotForAmount := strings.Split(amountEntry.Text, ".")
	// if amount is a float
	amountLabelBox := fyne.NewContainerWithLayout(layouts.NewHBox(0))
	if len(trailingDotForAmount) > 1 && len(trailingDotForAmount[1]) > 2 {
		trailingAmountLabel := canvas.NewText(trailingDotForAmount[1][2:]+" DCR", color.Black)
		trailingAmountLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
		trailingAmountLabel.TextSize = 15

		leadingAmountLabel := canvas.NewText(trailingDotForAmount[0]+"."+trailingDotForAmount[1][:2], color.Black)
		leadingAmountLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
		leadingAmountLabel.TextSize = 20

		amountLabelBox.AddObject(leadingAmountLabel)
		amountLabelBox.AddObject(trailingAmountLabel)

	} else {
		amountLabel := canvas.NewText(amountEntry.Text, color.Black)
		amountLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
		amountLabel.TextSize = 20

		DCRLabel := canvas.NewText("DCR", color.Black)
		DCRLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
		DCRLabel.TextSize = 15

		amountLabelBox.Layout = layouts.NewHBox(5)
		amountLabelBox.AddObject(amountLabel)
		amountLabelBox.AddObject(DCRLabel)
	}

	toDestination := "To destination address"
	destinationAddress := destinationAddressEntry.Text

	if sendingToSelf {
		toDestination = "To self"
		destinationAddress = selfSendingSelectedAccountText + " (" + sendingToSelfSelectedWalletName + ")"
	}

	sendButton := widgets.NewButton(color.RGBA{41, 112, 255, 255}, "Send "+amountEntry.Text+" DCR", func() {
		errorLabel := canvas.NewText("Wrong spending password. Please try again.", color.RGBA{237, 109, 71, 255})
		errorLabel.Alignment = fyne.TextAlignCenter
		errorLabel.TextSize = 12
		errorLabel.Hide()

		var confirmButton *widgets.Button

		walletPassword := widget.NewPasswordEntry()
		walletPassword.SetPlaceHolder("Spending Password")
		walletPassword.OnChanged = func(value string) {
			if value == "" {
				confirmButton.Disable()
			} else if confirmButton.Disabled() {
				confirmButton.Enable()
			}
		}

		var sendingPasswordPopup *widget.PopUp
		var popupContent *widget.Box

		cancelLabel := canvas.NewText("Cancel", color.RGBA{41, 112, 255, 255})
		cancelLabel.TextStyle.Bold = true

		cancelButton := widgets.NewClickableBox(widget.NewHBox(cancelLabel), func() {
			sendingPasswordPopup.Hide()
			confirmationPagePopup.Show()
		})

		confirmButton = widgets.NewButton(color.RGBA{41, 112, 255, 255}, "Confirm", func() {
			confirmButton.Disable()
			cancelButton.Disable()

			_, err := transactionAuthor.Broadcast([]byte(walletPassword.Text))
			if err != nil {
				// do not exit password popup on invalid passphrase
				if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
					errorLabel.Show()
					// this is an hack as selective refresh to errorLabel doesn't work
					popupContent.Refresh()
					confirmButton.Enable()
					cancelButton.Disable()
				} else {
					log.Println(err)
					errorLabelContainer.Show()
					sendingPasswordPopup.Hide()
					confirmationPagePopup.Show()
				}
				return
			}

			destinationAddressEntry.SetText("")
			amountEntry.SetText("")

			showSuccess.Container.Show()
			contents.Refresh()

			sendingPasswordPopup.Hide()

			time.AfterFunc(time.Second*5, func() {
				showSuccess.Container.Hide()
				contents.Refresh()
			})
		})
		confirmButton.SetMinSize(fyne.NewSize(91, 40))
		confirmButton.Disable()

		var passwordConceal *widgets.ImageButton
		passwordConceal = widgets.NewImageButton(reveal, nil, func() {
			if walletPassword.Password {
				passwordConceal.SetIcon(conceal)
				walletPassword.Password = false
			} else {
				passwordConceal.SetIcon(reveal)
				walletPassword.Password = true
			}
			// reveal texts
			walletPassword.SetText(walletPassword.Text)
		})

		popupContent = widget.NewHBox(
			widgets.NewHSpacer(24),
			widget.NewVBox(
				widgets.NewVSpacer(24),
				widget.NewLabelWithStyle("Confirm to send", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widgets.NewVSpacer(40),
				fyne.NewContainerWithLayout(layouts.NewPasswordLayout(fyne.NewSize(312, walletPassword.MinSize().Height)), walletPassword, passwordConceal),
				errorLabel,
				widgets.NewVSpacer(20),
				widget.NewHBox(layout.NewSpacer(), cancelButton, widgets.NewHSpacer(24), confirmButton.Container),
				widgets.NewVSpacer(24),
			),
			widgets.NewHSpacer(24),
		)

		sendingPasswordPopup = widget.NewModalPopUp(popupContent, window.Canvas())
	})

	sendButton.SetMinSize(fyne.NewSize(312, 56))
	sendButton.SetTextSize(18)

	confirmationPageContent := widget.NewVBox(
		widgets.NewVSpacer(18),
		accountSelectionPopupHeader,
		widgets.NewVSpacer(18),
		canvas.NewLine(color.Black),
		widgets.NewVSpacer(8),
		widget.NewHBox(layout.NewSpacer(), errorLabelContainer, layout.NewSpacer()),
		widgets.NewVSpacer(16),
		widget.NewHBox(layout.NewSpacer(), widget.NewLabel("Sending from"), sendingSelectedWalletLabel, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), amountLabelBox, layout.NewSpacer()),
		widgets.NewVSpacer(10),
		widget.NewHBox(layout.NewSpacer(), widget.NewIcon(downArrow), layout.NewSpacer()),
		widgets.NewVSpacer(10),
		widget.NewLabelWithStyle(toDestination, fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}),
		widget.NewLabelWithStyle(destinationAddress, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widgets.NewVSpacer(8),
		canvas.NewLine(color.RGBA{230, 234, 237, 255}),
		widget.NewHBox(widget.NewLabel("Transaction fee"),
			layout.NewSpacer(), widget.NewLabelWithStyle(transactionFeeText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		canvas.NewLine(color.RGBA{230, 234, 237, 255}),
		widget.NewHBox(widget.NewLabel("Total cost"),
			layout.NewSpacer(), widget.NewLabelWithStyle(totalCostText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		widget.NewHBox(widget.NewLabel("Balance after send"),
			layout.NewSpacer(), widget.NewLabelWithStyle(balanceAfterSendText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		canvas.NewLine(color.RGBA{230, 234, 237, 255}),
		widget.NewHBox(layout.NewSpacer(),
			widget.NewIcon(alert), widget.NewLabelWithStyle("Your DCR will be sent and CANNOT be undone.", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer()),
		sendButton.Container,
		widgets.NewVSpacer(18),
	)

	confirmationPagePopup = widget.NewModalPopUp(
		widget.NewHBox(widgets.NewHSpacer(16), confirmationPageContent, widgets.NewHSpacer(16)),
		window.Canvas())

	confirmationPagePopup.Show()
}

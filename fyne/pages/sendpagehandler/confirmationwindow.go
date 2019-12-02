package sendpagehandler

import (
	"fmt"
	"image/color"
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

		onConfirm := func(password string) error {
			_, err := transactionAuthor.Broadcast([]byte(password))
			return err
		}

		onCancel := func() {
			confirmationPagePopup.Show()
		}

		onError := func() {
			errorLabelContainer.Show()
			confirmationPagePopup.Show()
		}

		extraCalls := func() {
			destinationAddressEntry.SetText("")
			amountEntry.SetText("")

			showSuccess.Container.Show()
			contents.Refresh()

			time.AfterFunc(time.Second*5, func() {
				showSuccess.Container.Hide()
				contents.Refresh()
			})
		}

		PasswordPopUp(onConfirm, onCancel, onError, extraCalls, conceal, reveal, window)
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

package sendpagehandler

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"github.com/raedahgroup/godcr/fyne/assets"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/layouts"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

// todo: review confirmation window and check required parameters
func (sendPage *SendPageObjects) confirmationWindow() error {

	icons, err := assets.GetIcons(assets.CollapseDropdown, assets.ExpandDropdown, assets.DownArrow, assets.Alert, assets.Reveal)
	if err != nil {
		return errors.New("Unable to load icons")
	}

	var confirmationPagePopup *widget.PopUp

	confirmLabel := canvas.NewText("Confirm to send", color.Black)
	confirmLabel.TextStyle.Bold = true
	confirmLabel.TextSize = 20

	// check if truly it shows without passing widget renderer
	errorLabelContainer := widgets.NewBorderedText(failedToSend, fyne.NewSize(20, 16), color.RGBA{237, 109, 71, 255})
	errorLabelContainer.Container.Hide()

	accountSelectionPopupHeader := widget.NewHBox(
		widgets.NewImageButton(theme.CancelIcon(), nil, func() { confirmationPagePopup.Hide() }),
		widgets.NewHSpacer(9),
		confirmLabel,
		widgets.NewHSpacer(170),
	)
	sendingSelectedWalletLabel := widget.NewLabelWithStyle(fmt.Sprintf("%s (%s)",
		sendPage.Sending.SelectedAccountLabel.Text, sendPage.Sending.selectedWalletLabel.Text), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	trailingDotForAmount := strings.Split(sendPage.amountEntry.Text, ".")
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
		amountLabel := canvas.NewText(sendPage.amountEntry.Text, color.Black)
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
	destination := sendPage.destinationAddressEntry.Text
	destinationAddress := sendPage.destinationAddressEntry.Text

	if sendPage.destinationAddressEntry.Hidden {
		toDestination = "To self"
		destination = sendPage.SelfSending.SelectedAccountLabel.Text + " (" + sendPage.SelfSending.selectedWalletLabel.Text + ")"

		accNo, err := sendPage.SelfSending.selectedWallet.AccountNumber(sendPage.SelfSending.SelectedAccountLabel.Text)
		if err != nil {
			sendPage.showErrorLabel("Could not retrieve account number")
			return err
		}

		destinationAddress, err = sendPage.SelfSending.selectedWallet.CurrentAddress(int32(accNo))
		if err != nil {
			sendPage.showErrorLabel("Could not generate address to self send")
			return err
		}
	}

	sendButton := widgets.NewButton(color.RGBA{41, 112, 255, 255}, "Send "+sendPage.amountEntry.Text+" DCR", func() {
		onConfirm := func(password string) error {
			amountInFloat, err := strconv.ParseFloat(sendPage.amountEntry.Text, 64)
			if err != nil {
				sendPage.showErrorLabel("Could not parse float")
				return err
			}

			transactionAuthor, _ := sendPage.initTxAuthorAndGetAmountInWalletAccount(amountInFloat, destinationAddress)
			_, err = transactionAuthor.Broadcast([]byte(password))
			return err
		}

		onCancel := func() {
			confirmationPagePopup.Show()
		}

		onError := func() {
			errorLabelContainer.Container.Show()
			confirmationPagePopup.Show()
		}

		extraCalls := func() {
			sendPage.destinationAddressEntry.SetText("")
			sendPage.amountEntry.SetText("")

			sendPage.successLabel.Container.Show()
			sendPage.SendPageContents.Refresh()

			time.AfterFunc(time.Second*5, func() {
				sendPage.successLabel.Container.Hide()
				sendPage.SendPageContents.Refresh()
			})
		}

		passwordPopUp := PasswordPopUpObjects{
			onConfirm, onCancel, onError, extraCalls, sendPage.Window,
		}

		err := passwordPopUp.PasswordPopUp()
		if err != nil {

		}
	})

	sendButton.SetMinSize(fyne.NewSize(312, 56))
	sendButton.SetTextSize(18)

	confirmationPageContent := widget.NewVBox(
		widgets.NewVSpacer(18),
		accountSelectionPopupHeader,
		widgets.NewVSpacer(18),
		canvas.NewLine(color.Black),
		widgets.NewVSpacer(8),
		widget.NewHBox(layout.NewSpacer(), errorLabelContainer.Container, layout.NewSpacer()),
		widgets.NewVSpacer(16),
		widget.NewHBox(layout.NewSpacer(), widget.NewLabel("Sending from"), sendingSelectedWalletLabel, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), amountLabelBox, layout.NewSpacer()),
		widgets.NewVSpacer(10),
		widget.NewHBox(layout.NewSpacer(), widget.NewIcon(icons[assets.DownArrow]), layout.NewSpacer()),
		widgets.NewVSpacer(10),
		widgets.NewTextWithStyle(toDestination, color.RGBA{89, 109, 129, 255}, fyne.TextStyle{Bold: true}, fyne.TextAlignCenter, 14),
		widget.NewLabelWithStyle(destination, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widgets.NewVSpacer(8),
		canvas.NewLine(color.RGBA{230, 234, 237, 255}),
		widget.NewHBox(canvas.NewText("Transaction fee", color.RGBA{89, 109, 129, 255}),
			layout.NewSpacer(), widget.NewLabelWithStyle(sendPage.transactionFeeLabel.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		canvas.NewLine(color.RGBA{230, 234, 237, 255}),
		widget.NewHBox(canvas.NewText("Total cost", color.RGBA{89, 109, 129, 255}),
			layout.NewSpacer(), widget.NewLabelWithStyle(sendPage.totalCostLabel.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		widget.NewHBox(canvas.NewText("Balance after send", color.RGBA{89, 109, 129, 255}),
			layout.NewSpacer(), widget.NewLabelWithStyle(sendPage.balanceAfterSendLabel.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		canvas.NewLine(color.RGBA{230, 234, 237, 255}),
		widget.NewHBox(layout.NewSpacer(),
			widget.NewIcon(icons[assets.Alert]), widget.NewLabelWithStyle(sendingDcrWarning, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer()),
		sendButton.Container,
		widgets.NewVSpacer(18),
	)

	confirmationPagePopup = widget.NewModalPopUp(
		widget.NewHBox(widgets.NewHSpacer(16), confirmationPageContent, widgets.NewHSpacer(16)),
		sendPage.Window.Canvas())

	confirmationPagePopup.Show()

	return nil
}

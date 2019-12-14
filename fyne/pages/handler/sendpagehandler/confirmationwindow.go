package sendpagehandler

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/layouts"
	"github.com/raedahgroup/godcr/fyne/pages/handler/constantvalues"
	"github.com/raedahgroup/godcr/fyne/pages/handler/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

// todo: review confirmation window and check required parameters
func (sendPage *SendPageObjects) confirmationWindow() error {

	icons, err := assets.GetIcons(assets.CollapseDropdown, assets.ExpandDropdown, assets.DownArrow, assets.Alert, assets.Reveal)
	if err != nil {
		return errors.New(constantvalues.ConfirmationWindowIconsErr)
	}

	var confirmationPagePopup *widget.PopUp

	confirmLabel := widgets.NewTextWithStyle(constantvalues.ConfirmToSend, color.Black, fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, 20)

	// check if truly it shows without passing widget renderer
	errorLabelContainer := widgets.NewBorderedText(constantvalues.FailedToSend, fyne.NewSize(20, 16), color.RGBA{237, 109, 71, 255})
	errorLabelContainer.Container.Hide()

	accountSelectionPopupHeader := widget.NewHBox(
		widgets.NewImageButton(theme.CancelIcon(), nil, func() { confirmationPagePopup.Hide() }),
		widgets.NewHSpacer(9),
		confirmLabel,
		widgets.NewHSpacer(170),
	)
	sendingSelectedWalletLabel := widget.NewLabelWithStyle(fmt.Sprintf("%s (%s)",
		sendPage.Sending.SelectedAccountLabel.Text, sendPage.Sending.SelectedWalletLabel.Text), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	trailingDotForAmount := strings.Split(sendPage.amountEntry.Text, ".")
	// if amount is a float
	amountLabelBox := fyne.NewContainerWithLayout(layouts.NewHBox(0))
	if len(trailingDotForAmount) > 1 && len(trailingDotForAmount[1]) > 2 {
		trailingAmountLabel := widgets.NewTextWithStyle(fmt.Sprintf("%s %s", trailingDotForAmount[1][2:], constantvalues.DCR), color.Black, fyne.TextStyle{Bold: true, Monospace: true}, fyne.TextAlignLeading, 15)
		leadingAmountLabel := widgets.NewTextWithStyle(trailingDotForAmount[0]+"."+trailingDotForAmount[1][:2], color.Black, fyne.TextStyle{Bold: true, Monospace: true}, fyne.TextAlignLeading, 20)

		amountLabelBox.AddObject(leadingAmountLabel)
		amountLabelBox.AddObject(trailingAmountLabel)

	} else {
		amountLabel := widgets.NewTextWithStyle(sendPage.amountEntry.Text, color.Black, fyne.TextStyle{Bold: true, Monospace: true}, fyne.TextAlignLeading, 20)

		DCRLabel := widgets.NewTextWithStyle(constantvalues.DCR, color.Black, fyne.TextStyle{Bold: true, Monospace: true}, fyne.TextAlignLeading, 15)

		amountLabelBox.Layout = layouts.NewHBox(5)
		amountLabelBox.AddObject(amountLabel)
		amountLabelBox.AddObject(DCRLabel)
	}

	toDestination := constantvalues.ToDesinationAddress
	destination := sendPage.destinationAddressEntry.Text
	destinationAddress := sendPage.destinationAddressEntry.Text

	if sendPage.destinationAddressEntry.Hidden {
		toDestination = constantvalues.ToSelf
		destination = sendPage.SelfSending.SelectedAccountLabel.Text + " (" + sendPage.SelfSending.SelectedWalletLabel.Text + ")"

		accNo, err := sendPage.SelfSending.SelectedWallet.AccountNumber(sendPage.SelfSending.SelectedAccountLabel.Text)
		if err != nil {
			sendPage.showErrorLabel(constantvalues.AccountNumberErr)
			return err
		}

		destinationAddress, err = sendPage.SelfSending.SelectedWallet.CurrentAddress(int32(accNo))
		if err != nil {
			sendPage.showErrorLabel(constantvalues.GettingAddressToSelfSendErr)
			return err
		}
	}

	sendButton := widgets.NewButton(color.RGBA{41, 112, 255, 255}, fmt.Sprintf("%s %s %s", constantvalues.Send, sendPage.amountEntry.Text, constantvalues.DCR), func() {
		onConfirm := func(password string) error {
			amountInFloat, err := strconv.ParseFloat(sendPage.amountEntry.Text, 64)
			if err != nil {
				sendPage.showErrorLabel(constantvalues.ParseFloatErr)
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

		passwordPopUp := multipagecomponents.PasswordPopUpObjects{
			onConfirm, onCancel, onError, extraCalls, sendPage.Window,
		}

		passwordPopUp.PasswordPopUp()
	})

	sendButton.SetMinSize(fyne.NewSize(312, 56))
	sendButton.SetTextSize(18)

	confirmationPageContent := widget.NewVBox(
		widgets.NewVSpacer(10),
		accountSelectionPopupHeader,
		widgets.NewVSpacer(10),
		canvas.NewLine(color.Black),
		widgets.NewVSpacer(8),
		widget.NewHBox(layout.NewSpacer(), errorLabelContainer.Container, layout.NewSpacer()),
		widgets.NewVSpacer(16),
		widget.NewHBox(layout.NewSpacer(), widget.NewLabel(constantvalues.SendingFrom), sendingSelectedWalletLabel, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), amountLabelBox, layout.NewSpacer()),
		widgets.NewVSpacer(10),
		widget.NewHBox(layout.NewSpacer(), widget.NewIcon(icons[assets.DownArrow]), layout.NewSpacer()),
		widgets.NewVSpacer(10),
		widgets.NewTextWithStyle(toDestination, color.RGBA{89, 109, 129, 255}, fyne.TextStyle{Bold: true}, fyne.TextAlignCenter, 14),
		widget.NewLabelWithStyle(destination, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widgets.NewVSpacer(8),
		canvas.NewLine(color.RGBA{230, 234, 237, 255}),
		widgets.NewVSpacer(8),
		widget.NewHBox(canvas.NewText(constantvalues.TransactionFee, color.RGBA{89, 109, 129, 255}),
			layout.NewSpacer(), widgets.NewTextWithStyle(sendPage.transactionFeeLabel.Text, color.Black, fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, constantvalues.DefaultTextSize)),
		widgets.NewVSpacer(8),
		canvas.NewLine(color.RGBA{230, 234, 237, 255}),
		widgets.NewVSpacer(8),
		widget.NewHBox(canvas.NewText(constantvalues.TotalCost, color.RGBA{89, 109, 129, 255}),
			layout.NewSpacer(), widgets.NewTextWithStyle(sendPage.totalCostLabel.Text, color.Black, fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, constantvalues.DefaultTextSize)),
		widgets.NewVSpacer(4),
		widget.NewHBox(canvas.NewText(constantvalues.BalanceAfterSend, color.RGBA{89, 109, 129, 255}),
			layout.NewSpacer(), widgets.NewTextWithStyle(sendPage.balanceAfterSendLabel.Text, color.Black, fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, constantvalues.DefaultTextSize)),
		widgets.NewVSpacer(8),
		canvas.NewLine(color.RGBA{230, 234, 237, 255}),
		widgets.NewVSpacer(8),
		widget.NewHBox(layout.NewSpacer(),
			widget.NewIcon(icons[assets.Alert]), widgets.NewTextWithStyle(constantvalues.SendingDcrWarning, color.Black, fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, constantvalues.DefaultTextSize), layout.NewSpacer()),
		widgets.NewVSpacer(8),
		sendButton.Container,
		widgets.NewVSpacer(18),
	)

	confirmationPagePopup = widget.NewModalPopUp(
		widget.NewHBox(widgets.NewHSpacer(16), confirmationPageContent, widgets.NewHSpacer(16)),
		sendPage.Window.Canvas())

	confirmationPagePopup.Show()

	return nil
}

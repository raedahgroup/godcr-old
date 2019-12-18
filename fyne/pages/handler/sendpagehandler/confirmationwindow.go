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
	"github.com/raedahgroup/godcr/fyne/pages/handler/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) confirmationWindow() error {
	icons, err := assets.GetIcons(assets.CollapseDropdown, assets.ExpandDropdown, assets.DownArrow, assets.Alert, assets.Reveal)
	if err != nil {
		return errors.New(values.ConfirmationWindowIconsErr)
	}

	var confirmationPagePopup *widget.PopUp

	confirmLabel := widgets.NewTextWithStyle(values.ConfirmToSend, color.Black, fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, 20)

	errorLabelContainer := widgets.NewBorderedText(values.FailedToSend, fyne.NewSize(20, 16), values.ErrorColor)
	errorLabelContainer.Container.Hide()

	accountSelectionPopupHeader := widget.NewHBox(
		widgets.NewImageButton(theme.CancelIcon(), nil, func() { confirmationPagePopup.Hide() }),
		widgets.NewHSpacer(values.SpacerSize10),
		confirmLabel,
		widgets.NewHSpacer(values.SpacerSize170),
	)
	sendingSelectedWalletLabel := widget.NewLabelWithStyle(fmt.Sprintf("%s (%s)",
		sendPage.Sending.SelectedAccountLabel.Text, sendPage.Sending.SelectedWalletLabel.Text), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	trailingDotForAmount := strings.Split(sendPage.amountEntry.Text, ".")
	// if amount is a float
	amountLabelBox := fyne.NewContainerWithLayout(layouts.NewHBox(0))
	if len(trailingDotForAmount) > 1 && len(trailingDotForAmount[1]) > 2 {
		trailingAmountLabel := widgets.NewTextWithStyle(fmt.Sprintf("%s %s", trailingDotForAmount[1][2:], values.DCR),
			values.DefaultTextColor, fyne.TextStyle{Bold: true, Monospace: true}, fyne.TextAlignLeading, values.TextSize16)
		leadingAmountLabel := widgets.NewTextWithStyle(trailingDotForAmount[0]+"."+trailingDotForAmount[1][:2],
			values.DefaultTextColor, fyne.TextStyle{Bold: true, Monospace: true}, fyne.TextAlignLeading, values.TextSize20)

		amountLabelBox.AddObject(leadingAmountLabel)
		amountLabelBox.AddObject(trailingAmountLabel)

	} else {
		amountLabel := widgets.NewTextWithStyle(sendPage.amountEntry.Text, color.Black, fyne.TextStyle{Bold: true, Monospace: true}, fyne.TextAlignLeading, values.SpacerSize20)

		DCRLabel := widgets.NewTextWithStyle(values.DCR, color.Black, fyne.TextStyle{Bold: true, Monospace: true}, fyne.TextAlignLeading, values.SpacerSize16)

		amountLabelBox.Layout = layouts.NewHBox(values.SpacerSize4)
		amountLabelBox.AddObject(amountLabel)
		amountLabelBox.AddObject(DCRLabel)
	}

	toDestination := values.ToDesinationAddress
	destination := sendPage.destinationAddressEntry.Text
	destinationAddress := sendPage.destinationAddressEntry.Text

	if sendPage.destinationAddressEntry.Hidden {
		toDestination = values.ToSelf
		destination = sendPage.SelfSending.SelectedAccountLabel.Text + " (" + sendPage.SelfSending.SelectedWalletLabel.Text + ")"

		destinationAddress, err = sendPage.SelfSending.SelectedWallet.CurrentAddress(int32(*sendPage.SelfSending.SendingSelectedAccountID))
		if err != nil {
			sendPage.showErrorLabel(values.GettingAddressToSelfSendErr)
			return err
		}
	}

	sendButton := widgets.NewButton(values.Blue, fmt.Sprintf("%s %s %s", values.Send, sendPage.amountEntry.Text, values.DCR), func() {
		onConfirm := func(password string) error {
			amountInFloat, err := strconv.ParseFloat(sendPage.amountEntry.Text, 64)
			if err != nil {
				sendPage.showErrorLabel(values.ParseFloatErr)
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
	sendButton.SetTextSize(values.ConfirmationButtonTextSize)

	confirmationPageContent := widget.NewVBox(
		widgets.NewVSpacer(values.SpacerSize10),
		accountSelectionPopupHeader,
		widgets.NewVSpacer(values.SpacerSize10),
		canvas.NewLine(color.Black),
		widgets.NewVSpacer(values.SpacerSize8),
		widget.NewHBox(layout.NewSpacer(), errorLabelContainer.Container, layout.NewSpacer()),
		widgets.NewVSpacer(values.SpacerSize16),

		widget.NewHBox(layout.NewSpacer(), fyne.NewContainerWithLayout(layouts.NewHBox(values.NilSpacer),
			widget.NewLabel(values.SendingFrom), sendingSelectedWalletLabel), layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), amountLabelBox, layout.NewSpacer()),
		widgets.NewVSpacer(values.SpacerSize10),
		widget.NewHBox(layout.NewSpacer(), widget.NewIcon(icons[assets.DownArrow]), layout.NewSpacer()),
		widgets.NewVSpacer(values.SpacerSize10),
		widgets.NewTextWithStyle(toDestination, values.TransactionInfoColor, fyne.TextStyle{}, fyne.TextAlignCenter, values.DefaultTextSize),
		widgets.NewTextWithStyle(destination, values.DefaultTextColor, fyne.TextStyle{}, fyne.TextAlignCenter, values.DefaultTextSize),
		widgets.NewVSpacer(values.SpacerSize8),

		canvas.NewLine(values.StrippedLineColor),
		widgets.NewVSpacer(values.SpacerSize8),
		widget.NewHBox(canvas.NewText(values.TransactionFee, values.TransactionInfoColor),
			layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%s %s", sendPage.transactionFeeLabel.Text, values.DCR), values.DefaultTextColor)),
		widgets.NewVSpacer(values.SpacerSize8),
		canvas.NewLine(values.StrippedLineColor),
		widgets.NewVSpacer(values.SpacerSize8),
		widget.NewHBox(canvas.NewText(values.TotalCost, values.TransactionInfoColor),
			layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%s %s", sendPage.totalCostLabel.Text, values.DCR), values.DefaultTextColor)),
		widgets.NewVSpacer(values.SpacerSize4),
		widget.NewHBox(canvas.NewText(values.BalanceAfterSend, values.TransactionInfoColor),
			layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%s %s", sendPage.balanceAfterSendLabel.Text, values.DCR), values.DefaultTextColor)),
		widgets.NewVSpacer(values.SpacerSize8),
		canvas.NewLine(values.StrippedLineColor),
		widgets.NewVSpacer(values.SpacerSize8),
		widget.NewHBox(layout.NewSpacer(),
			widget.NewIcon(icons[assets.Alert]), canvas.NewText(values.SendingDcrWarning, values.DefaultTextColor), layout.NewSpacer()),
		widgets.NewVSpacer(values.SpacerSize8),
		sendButton.Container,
		widgets.NewVSpacer(values.SpacerSize18),
	)

	confirmationPagePopup = widget.NewModalPopUp(
		widget.NewHBox(widgets.NewHSpacer(values.SpacerSize16), confirmationPageContent, widgets.NewHSpacer(values.SpacerSize16)),
		sendPage.Window.Canvas())

	confirmationPagePopup.Show()

	return nil
}

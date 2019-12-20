package sendpagehandler

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/fyne/widgets"

	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
)

// sendingToDestinationComponents constitutes all components that composes sending coins to users or self.
func (sendPage *SendPageObjects) initToDestinationComponents() error {
	fromLabel := canvas.NewText("To", values.DarkerBlueGrayTextColor)

	accountBox, err := sendPage.SelfSending.CreateAccountSelector("Sending account")
	if err != nil {
		return err
	}
	accountBox.Hide()

	sendPage.destinationAddressEntryComponent()

	sendToAccountLabel := canvas.NewText(values.SwitchToSendToAccount, values.Blue)
	sendToAccountLabel.TextSize = values.TextSize12

	destinationAddressContainer := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(
		fyne.NewSize(widget.NewLabel(values.TestAddress).MinSize().Width, sendPage.destinationAddressEntry.MinSize().Height)), sendPage.destinationAddressEntry)

	spacer := widgets.NewVSpacer(values.SpacerSize10)

	var container *fyne.Container
	switchingComponentButton := widgets.NewClickableBox(widget.NewVBox(sendToAccountLabel), func() {
		sendPage.SendPageContents.Refresh()
		if accountBox.Hidden {
			sendToAccountLabel.Text = values.SwitchToSendToAddress
			accountBox.Show()
			sendPage.destinationAddressEntry.Hide()
			destinationAddressContainer.Hide()
			sendPage.destinationAddressErrorLabel.Hide()
			spacer.Hide()

		} else {
			sendToAccountLabel.Text = values.SwitchToSendToAccount
			sendPage.destinationAddressEntry.Show()
			destinationAddressContainer.Show()
			accountBox.Hide()
			spacer.Show()

			sendPage.destinationAddressEntry.OnChanged(sendPage.destinationAddressEntry.Text)
		}

		sendPage.SendPageContents.Refresh()
		sendPage.initTxDetails(sendPage.amountEntry.Text)
		sendPage.SendPageContents.Refresh()
	})

	box := widget.NewVBox(
		widget.NewHBox(fromLabel, layout.NewSpacer(), switchingComponentButton, widgets.NewHSpacer(values.SpacerSize20)),
		accountBox,
		destinationAddressContainer,
		sendPage.destinationAddressErrorLabel,
		spacer)

	container = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(box.MinSize()), box)

	sendPage.SendPageContents.Append(container)

	return err
}

func (sendPage *SendPageObjects) destinationAddressEntryComponent() {
	sendPage.destinationAddressErrorLabel = canvas.NewText("", values.ErrorColor)
	sendPage.destinationAddressErrorLabel.TextSize = values.DefaultErrTextSize
	sendPage.destinationAddressErrorLabel.Hide()

	sendPage.destinationAddressEntry = widget.NewEntry()
	sendPage.destinationAddressEntry.SetPlaceHolder(values.DestinationAddressPlaceHolder)

	sendPage.destinationAddressEntry.OnChanged = func(address string) {
		if sendPage.destinationAddressEntry.Text == "" {
			sendPage.destinationAddressErrorLabel.Hide()
			sendPage.SendPageContents.Refresh()
			sendPage.initTxDetails(sendPage.amountEntry.Text)

			return
		}

		_, err := dcrutil.DecodeAddress(address)
		if err != nil {
			sendPage.destinationAddressErrorLabel.Text = values.InvalidAddress
			sendPage.SendPageContents.Refresh()
			sendPage.destinationAddressErrorLabel.Show()
			sendPage.SendPageContents.Refresh()
			setLabelText(values.NilAmount, sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel)
			setLabelText(values.ZeroByte, sendPage.transactionSizeLabel)
			sendPage.SendPageContents.Refresh()

		} else {
			sendPage.destinationAddressErrorLabel.Hide()
			sendPage.SendPageContents.Refresh()
		}

		if sendPage.amountEntry.Text != "" && sendPage.amountEntryErrorLabel.Hidden && sendPage.destinationAddressErrorLabel.Hidden {
			sendPage.nextButton.Enable()
		} else {
			sendPage.nextButton.Disable()
		}

		sendPage.SendPageContents.Refresh()
	}
}

func setLabelText(Text string, objects ...*canvas.Text) {
	for _, object := range objects {
		object.Text = Text
	}
}

func setLabelColor(textColor color.Color, objects ...*canvas.Text) {
	for _, object := range objects {
		object.Color = textColor
	}
}

func canvasTextRefresher(objects ...*canvas.Text) {
	for _, object := range objects {
		object.Refresh()
	}
}

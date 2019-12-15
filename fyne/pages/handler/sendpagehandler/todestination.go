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
	fromLabel := canvas.NewText("To", color.RGBA{61, 88, 115, 255})
	fromLabel.TextStyle.Bold = true

	accountBox, err := sendPage.SelfSending.CreateAccountSelector("Sending account")
	if err != nil {
		return err
	}
	accountBox.Hide()

	sendPage.destinationAddressEntryComponent()

	sendToAccountLabel := canvas.NewText(values.SwitchToSendToAccount, color.RGBA{R: 41, G: 112, B: 255, A: 255})
	sendToAccountLabel.TextSize = values.TextSize12

	destinationAddressContainer := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(
		fyne.NewSize(widget.NewLabel(values.TestAddress).MinSize().Width, sendPage.destinationAddressEntry.MinSize().Height)), sendPage.destinationAddressEntry)

	spacer := widgets.NewVSpacer(values.SpacerSize10)

	var container *fyne.Container
	switchingComponentButton := widgets.NewClickableBox(widget.NewVBox(sendToAccountLabel), func() {
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
		sendPage.amountEntry.OnChanged(sendPage.amountEntry.Text)
		sendPage.SendPageContents.Refresh()
	})

	box := widget.NewVBox(
		widget.NewHBox(fromLabel, layout.NewSpacer(), switchingComponentButton, widgets.NewHSpacer(20)),
		accountBox,
		destinationAddressContainer,
		sendPage.destinationAddressErrorLabel,
		spacer)

	container = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(box.MinSize()), box)

	sendPage.SendPageContents.Append(container)

	return err
}

func (sendPage *SendPageObjects) destinationAddressEntryComponent() {

	sendPage.destinationAddressErrorLabel = canvas.NewText("", color.RGBA{237, 109, 71, 255})
	sendPage.destinationAddressErrorLabel.TextSize = values.DefaultErrTextSize
	sendPage.destinationAddressErrorLabel.Hide()

	sendPage.destinationAddressEntry = widget.NewEntry()
	sendPage.destinationAddressEntry.SetPlaceHolder(values.DestinationAddressPlaceHolder)

	sendPage.destinationAddressEntry.OnChanged = func(address string) {
		if sendPage.destinationAddressEntry.Text == "" {
			sendPage.destinationAddressErrorLabel.Hide()
			sendPage.SendPageContents.Refresh()
			return
		}

		_, err := dcrutil.DecodeAddress(address)
		if err != nil {
			sendPage.destinationAddressErrorLabel.Text = values.InvalidAddress
			sendPage.destinationAddressErrorLabel.Show()
			setLabelText(values.NilAmount, sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel)
			setLabelText(values.ZeroByte, sendPage.transactionSizeLabel)

		} else {
			sendPage.destinationAddressErrorLabel.Hide()
		}

		if sendPage.amountEntry.Text != "" && sendPage.amountEntryErrorLabel.Hidden && sendPage.destinationAddressErrorLabel.Hidden {
			sendPage.nextButton.Enable()
		} else {
			sendPage.nextButton.Disable()
		}

		sendPage.SendPageContents.Refresh()
	}
}

func setLabelText(Text string, objects ...*widget.Label) {
	for _, object := range objects {
		object.SetText(Text)
		object.Refresh()
	}
}

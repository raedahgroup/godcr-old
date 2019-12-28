package send

import (
	"fyne.io/fyne"

	"github.com/raedahgroup/godcr/fyne/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initNextButton() {
	sendPage.nextButton = widgets.NewButton(values.Blue, "Next", func() {
		if sendPage.MultiWallet.ConnectedPeers() <= 0 {
			sendPage.showErrorLabel(values.NotConnectedErr)
			return
		}

		if sendPage.SelfSending.SelectedWallet == nil {
			sendPage.showErrorLabel(values.SelectedWalletInvalidErr)
			return
		}

		err := sendPage.confirmationWindow()
		if err != nil {
			sendPage.showErrorLabel(values.ConfirmationWindowErr)
		}
	})

	sendPage.nextButton.SetTextSize(values.ConfirmationButtonTextSize)
	sendPage.nextButton.SetMinSize(sendPage.nextButton.MinSize().Add(fyne.NewSize(0, 20)))
	sendPage.nextButton.Disable()
	sendPage.SendPageContents.Append(sendPage.nextButton.Container)
}

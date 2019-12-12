package sendpagehandler

import (
	"image/color"

	"fyne.io/fyne"

	"github.com/raedahgroup/godcr/fyne/pages/handler/constantvalues"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initNextButton() {

	sendPage.nextButton = widgets.NewButton(color.RGBA{41, 112, 255, 255}, "Next", func() {
		if sendPage.MultiWallet.ConnectedPeers() <= 0 {
			sendPage.showErrorLabel(constantvalues.NotConnectedErr)
			return
		}

		if sendPage.SelfSending.SelectedWallet == nil {
			sendPage.showErrorLabel(constantvalues.SelectedWalletInvalidErr)
			return
		}

		err := sendPage.confirmationWindow()
		if err != nil {
			sendPage.showErrorLabel(constantvalues.ConfirmationWindowErr)
		}
	})

	sendPage.nextButton.SetMinSize(sendPage.nextButton.MinSize().Add(fyne.NewSize(0, 20)))
	sendPage.nextButton.Disable()
	sendPage.SendPageContents.Append(sendPage.nextButton.Container)
}

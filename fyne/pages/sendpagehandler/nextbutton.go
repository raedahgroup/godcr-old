package sendpagehandler

import (
	"image/color"

	"fyne.io/fyne"

	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initNextButton() {

	sendPage.nextButton = widgets.NewButton(color.RGBA{41, 112, 255, 255}, "Next", func() {
		if sendPage.MultiWallet.ConnectedPeers() <= 0 {
			sendPage.showErrorLabel("Not Connected To Decred Network")
			return
		}

		if sendPage.SelfSending.selectedWallet == nil {
			sendPage.showErrorLabel("Selected self sending wallet is invalid")
			return
		}

		err := sendPage.confirmationWindow()
		if err != nil {
			sendPage.showErrorLabel("Could not view confirmation window")
		}
	})

	sendPage.nextButton.SetMinSize(sendPage.nextButton.MinSize().Add(fyne.NewSize(0, 20)))
	sendPage.nextButton.Disable()
	sendPage.SendPageContents.Append(sendPage.nextButton.Container)
}

package receivepagehandler

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/skip2/go-qrcode"

	"github.com/raedahgroup/godcr/fyne/pages/handler/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type ReceivePageObjects struct {
	Accounts    multipagecomponents.AccountSelectorStruct
	MultiWallet *dcrlibwallet.MultiWallet

	qrImage             *widget.Icon
	address             *canvas.Text
	addressCopiedLabel  *widgets.BorderedText
	errorLabel          *widgets.BorderedText
	ReceivePageContents *widget.Box

	Window fyne.Window
}

func (receivePage *ReceivePageObjects) InitReceivePage() error {
	receivePage.ReceivePageContents.Append(widgets.NewVSpacer(values.Padding))

	err := receivePage.initBaseObjects()
	if err != nil {
		return err
	}

	receivePage.ReceivePageContents.Append(widgets.NewVSpacer(values.SpacerSize10))

	receivePage.errorLabel = widgets.NewBorderedText("", fyne.NewSize(0, 0), color.RGBA{237, 109, 71, 255})
	receivePage.errorLabel.Container.Hide()

	receivePage.addressCopiedLabel = widgets.NewBorderedText("", fyne.NewSize(0, 0), color.RGBA{65, 190, 83, 255})
	receivePage.addressCopiedLabel.Container.Hide()

	receivePage.ReceivePageContents.Append(widget.NewHBox(layout.NewSpacer(), receivePage.errorLabel.Container, layout.NewSpacer()))

	err = receivePage.initAccountSelector()
	if err != nil {
		return err
	}

	receivePage.ReceivePageContents.Append(widget.NewHBox(layout.NewSpacer(), receivePage.addressCopiedLabel.Container, layout.NewSpacer()))

	receivePage.ReceivePageContents.Append(widgets.NewVSpacer(values.SpacerSize10))
	receivePage.initQrImageAndAddress()
	receivePage.initTapToCopyText()

	receivePage.ReceivePageContents.Append(widgets.NewVSpacer(values.Padding))

	return nil
}

func (receivePage *ReceivePageObjects) generateAddressAndQR(newAddress bool) {
	accNo, err := receivePage.Accounts.SelectedWallet.AccountNumber(receivePage.Accounts.SelectedAccountLabel.Text)
	if err != nil {
		receivePage.showInfoLabel(receivePage.errorLabel, values.AccountNumberErr)
		return
	}

	var addr string
	if newAddress {
		addr, err = receivePage.Accounts.SelectedWallet.NextAddress(int32(accNo))
		if err != nil {
			receivePage.showInfoLabel(receivePage.errorLabel, values.GettingAddress)
			return
		}
	} else {
		addr, err = receivePage.Accounts.SelectedWallet.CurrentAddress(int32(accNo))
		if err != nil {
			receivePage.showInfoLabel(receivePage.errorLabel, values.GettingAddress)
			return
		}
	}

	receivePage.address.Refresh()
	receivePage.address.Text = addr
	receivePage.address.Refresh()

	imgBytes, err := qrcode.Encode(addr, qrcode.High, 256)
	if err != nil {
		receivePage.showInfoLabel(receivePage.errorLabel, values.QrEncodeErr)
		return
	}

	receivePage.qrImage.SetResource(fyne.NewStaticResource("Text", imgBytes))

	receivePage.ReceivePageContents.Refresh()
}

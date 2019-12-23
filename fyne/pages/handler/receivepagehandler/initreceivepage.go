package receivepagehandler

import (
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
	borderedContent     *widgets.Box
	Window              fyne.Window
}

func (receivePage *ReceivePageObjects) InitReceivePage() error {
	receivePage.ReceivePageContents.Append(widgets.NewVSpacer(values.Padding))
	receivePage.borderedContent = widgets.NewVBox()

	err := receivePage.initBaseObjects()
	if err != nil {
		return err
	}

	receivePage.ReceivePageContents.Append(widgets.NewVSpacer(values.SpacerSize10))

	receivePage.errorLabel = widgets.NewBorderedText("", fyne.NewSize(0, 0), values.ErrorColor)
	receivePage.errorLabel.Container.Hide()

	receivePage.addressCopiedLabel = widgets.NewBorderedText("", fyne.NewSize(0, 0), values.Green)
	receivePage.addressCopiedLabel.Container.Hide()

	receivePage.borderedContent.Append(widgets.NewHBox(layout.NewSpacer(), receivePage.errorLabel.Container, layout.NewSpacer()))

	err = receivePage.initAccountSelector()
	if err != nil {
		return err
	}

	receivePage.borderedContent.Append(widgets.NewHBox(layout.NewSpacer(), receivePage.addressCopiedLabel.Container, layout.NewSpacer()))

	receivePage.borderedContent.Append(widgets.NewVSpacer(values.SpacerSize10))
	receivePage.initQrImageAndAddress()
	receivePage.initTapToCopyText()
	receivePage.borderedContent.Append(widgets.NewVSpacer(values.SpacerSize20))
	receivePage.ReceivePageContents.Append(receivePage.borderedContent)

	receivePage.ReceivePageContents.Append(widgets.NewVSpacer(values.Padding))

	return nil
}

func (receivePage *ReceivePageObjects) generateAddressAndQR(newAddress bool) {
	var addr string
	var err error
	if newAddress {
		addr, err = receivePage.Accounts.SelectedWallet.NextAddress(int32(*receivePage.Accounts.SendingSelectedAccountID))
		if err != nil {
			receivePage.showInfoLabel(receivePage.errorLabel, values.GettingAddress)
			return
		}
	} else {
		addr, err = receivePage.Accounts.SelectedWallet.CurrentAddress(int32(*receivePage.Accounts.SendingSelectedAccountID))
		if err != nil {
			receivePage.showInfoLabel(receivePage.errorLabel, values.GettingAddress)
			return
		}
	}

	receivePage.address.Refresh()
	receivePage.address.Text = addr
	receivePage.address.Refresh()

	qrCode, err := qrcode.New(addr, qrcode.Highest)
	qrCode.DisableBorder = true
	if err != nil {
		receivePage.showInfoLabel(receivePage.errorLabel, values.QrEncodeErr)
		return
	}

	bytes, err := qrCode.PNG(256)
	if err != nil {
		receivePage.showInfoLabel(receivePage.errorLabel, values.QrEncodeErr)
		return
	}

	receivePage.qrImage.SetResource(fyne.NewStaticResource("", bytes))

	receivePage.ReceivePageContents.Refresh()
}

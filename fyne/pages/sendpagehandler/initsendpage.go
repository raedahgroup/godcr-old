package sendpagehandler

import (
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/constantvalues"
	"github.com/raedahgroup/godcr/fyne/pages/multipagecomponents.go"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type SendPageObjects struct {
	errorLabel   *widgets.BorderedText
	successLabel *widgets.BorderedText

	destinationAddressEntry      *widget.Entry
	destinationAddressErrorLabel *canvas.Text

	amountEntry           *widget.Entry
	amountEntryErrorLabel *canvas.Text
	SpendableLabel        *canvas.Text

	transactionFeeLabel   *widget.Label
	transactionSizeLabel  *widget.Label
	totalCostLabel        *widget.Label
	balanceAfterSendLabel *widget.Label

	nextButton *widgets.Button

	Sending     multipagecomponents.AccountSelectorStruct
	SelfSending multipagecomponents.AccountSelectorStruct

	SendPageContents *widget.Box

	MultiWallet *dcrlibwallet.MultiWallet

	Window fyne.Window
}

func (sendPage *SendPageObjects) InitAllSendPageComponents() error {
	// add padding to the top
	sendPage.SendPageContents.Append(widgets.NewVSpacer(20))

	err := sendPage.initBaseObjects()
	if err != nil {
		return err
	}

	sendPage.SendPageContents.Append(widgets.NewVSpacer(10))

	sendPage.successLabel = widgets.NewBorderedText(constantvalues.SuccessText, fyne.NewSize(0, 0), color.RGBA{65, 190, 83, 255})
	sendPage.successLabel.Container.Hide()

	sendPage.errorLabel = widgets.NewBorderedText("", fyne.NewSize(0, 0), color.RGBA{237, 109, 71, 255})
	sendPage.errorLabel.Container.Hide()

	sendPage.SendPageContents.Append(widget.NewHBox(layout.NewSpacer(), sendPage.successLabel.Container, sendPage.errorLabel.Container, layout.NewSpacer()))

	err = sendPage.initFromAccountSelector()
	if err != nil {
		return err
	}

	sendPage.SendPageContents.Append(widgets.NewVSpacer(10))

	err = sendPage.initToDestinationComponents()
	if err != nil {
		return err
	}

	sendPage.SendPageContents.Append(widgets.NewVSpacer(20))

	sendPage.initAmountEntryComponents()

	sendPage.SendPageContents.Append(widgets.NewVSpacer(20))

	err = sendPage.initTransactionDetails()
	if err != nil {
		return err
	}

	sendPage.SendPageContents.Append(widgets.NewVSpacer(35))

	sendPage.initNextButton()
	return err
}

func (sendPage *SendPageObjects) showErrorLabel(err string) {
	sendPage.errorLabel.SetText(err)
	sendPage.errorLabel.SetPadding(fyne.NewSize(20, 8))
	sendPage.errorLabel.Container.Show()
	sendPage.SendPageContents.Refresh()

	time.AfterFunc(time.Second*5, func() {
		sendPage.errorLabel.Container.Hide()
		sendPage.SendPageContents.Refresh()
	})
}

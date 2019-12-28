package send

import (
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/values"
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

	transactionFeeLabel   *canvas.Text
	transactionSizeLabel  *canvas.Text
	totalCostLabel        *canvas.Text
	balanceAfterSendLabel *canvas.Text

	maxButton  *widgets.Button
	nextButton *widgets.Button

	Sending     multipagecomponents.AccountSelectorStruct
	SelfSending multipagecomponents.AccountSelectorStruct

	SendPageContents *widget.Box

	MultiWallet *dcrlibwallet.MultiWallet
	sendMax     bool

	Window fyne.Window
}

func (sendPage *SendPageObjects) InitAllSendPageComponents() error {
	sendPage.SendPageContents.Append(widgets.NewVSpacer(values.Padding)) // top padding

	err := sendPage.initBaseObjects()
	if err != nil {
		return err
	}

	sendPage.SendPageContents.Append(widgets.NewVSpacer(values.SpacerSize10))

	sendPage.successLabel = widgets.NewBorderedText(values.SuccessText, fyne.NewSize(0, 0), values.Green)
	sendPage.successLabel.Container.Hide()

	sendPage.errorLabel = widgets.NewBorderedText("", fyne.NewSize(0, 0), values.ErrorColor)
	sendPage.errorLabel.Container.Hide()

	sendPage.SendPageContents.Append(widget.NewHBox(layout.NewSpacer(), sendPage.successLabel.Container, sendPage.errorLabel.Container, layout.NewSpacer()))

	err = sendPage.initFromAccountSelector()
	if err != nil {
		return err
	}

	sendPage.SendPageContents.Append(widgets.NewVSpacer(values.SpacerSize10))

	err = sendPage.initToDestinationComponents()
	if err != nil {
		return err
	}

	sendPage.SendPageContents.Append(widgets.NewVSpacer(values.SpacerSize20))

	sendPage.initAmountEntryComponents()

	sendPage.SendPageContents.Append(widgets.NewVSpacer(values.SpacerSize36))

	err = sendPage.initTransactionDetails()
	if err != nil {
		return err
	}

	sendPage.SendPageContents.Append(widgets.NewVSpacer(values.SpacerSize30))

	sendPage.initNextButton()

	sendPage.SendPageContents.Append(widgets.NewVSpacer(values.BottomPadding))
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

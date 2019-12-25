package stakingpagehandler

import (
	"log"
	"regexp"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (stakingPage *StakingPageObjects) purchaseTickets() error {
	numberOfTicketsEntry := widget.NewEntry()
	numberOfTicketsEntry.PlaceHolder = "No. of purchase tickets"
	amountEntryExpression, err := regexp.Compile(values.AmountRegExp)
	if err != nil {
		log.Println(err)
	}

	numberOfTicketsEntry.OnChanged = func(value string) {
		if len(value) > 0 && !amountEntryExpression.MatchString(value) {
			if len(value) == 1 {
				numberOfTicketsEntry.SetText("")
			} else {
				//fix issue with crash on paste here
				value = value[:numberOfTicketsEntry.CursorColumn-1] + value[numberOfTicketsEntry.CursorColumn:]
				//todo: using setText, cursor column count doesnt increase or reduce. Create an issue on this
				numberOfTicketsEntry.CursorColumn--
				numberOfTicketsEntry.SetText(value)
			}

			return
		}
	}

	spendUnconfirmedCheck := widget.NewCheck("", nil)

	passphraseEntry := widget.NewPasswordEntry()
	passphraseEntry.PlaceHolder = "Enter spending Passphrase"

	accountBox, err := stakingPage.Accounts.CreateAccountSelector(values.ReceivingAccountLabel)
	if err != nil {
		return err
	}

	submitButton := widget.NewButton("Submit", func() {

	})

	purchaseForm := widget.NewForm()
	purchaseForm.Append("Source Account: ", accountBox)
	purchaseForm.Append("No. of tickets: ", numberOfTicketsEntry)
	purchaseForm.Append("Spend unconfirmed", spendUnconfirmedCheck)
	purchaseForm.Append("Passphrase", passphraseEntry)

	purchaseTicketsData := widget.NewVBox(
		purchaseForm,
		widgets.NewVSpacer(values.SpacerSize30),
		submitButton,
	)

	purchaseFormLayout := fyne.NewContainerWithLayout(layout.NewCenterLayout(), purchaseTicketsData)

	stakingPage.StakingPageContents.Append(purchaseFormLayout)
	return nil
}

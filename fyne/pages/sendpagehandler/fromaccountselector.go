package sendpagehandler

import (
	"image/color"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initFromAccountSelector() error {
	fromLabel := widgets.NewTextWithStyle("From", color.RGBA{61, 88, 115, 255}, fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, 14)

	sendPage.Sending.onAccountChange = sendPage.onAccountChange

	accountBox, err := sendPage.Sending.CreateAccountSelector("Sending account")
	if err != nil { // return err if icons in account selector dont load
		return err
	}

	box := widget.NewVBox(fromLabel, accountBox)

	sendPage.SendPageContents.Append(box)
	return err
}

func (sendPage *SendPageObjects) onAccountChange() {
	accountNumber, err := sendPage.Sending.selectedWallet.AccountNumber(sendPage.Sending.SelectedAccountLabel.Text)
	if err != nil {
		sendPage.showErrorLabel("Could not get accounts")
		log.Println("could not get accounts on account change, reason:", err.Error())
		return
	}

	balance, err := sendPage.Sending.selectedWallet.GetAccountBalance(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		sendPage.showErrorLabel("could not retrieve account balance")
		log.Println("could not retrieve account balance on account change, reason:", err.Error())
		return
	}
	sendPage.SpendableLabel.Text = "Spendable: " + dcrutil.Amount(balance.Spendable).String()
	sendPage.SpendableLabel.Refresh()

	sendPage.transactionFeeLabel.Refresh()
	sendPage.transactionSizeLabel.Refresh()
	sendPage.amountEntry.OnChanged(sendPage.amountEntry.Text)
}

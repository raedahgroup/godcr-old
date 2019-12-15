package sendpagehandler

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initFromAccountSelector() error {
	fromLabel := widgets.NewTextWithStyle(values.FromText, color.RGBA{61, 88, 115, 255}, fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, values.DefaultTextSize)

	sendPage.Sending.OnAccountChange = sendPage.onAccountChange

	accountBox, err := sendPage.Sending.CreateAccountSelector(values.FromAccountSelectorPopUpHeaderLabel)
	if err != nil { // return err if icons in account selector dont load
		return err
	}

	box := widget.NewVBox(fromLabel, accountBox)

	sendPage.SendPageContents.Append(box)
	return err
}

func (sendPage *SendPageObjects) onAccountChange() {
	accountNumber, err := sendPage.Sending.SelectedWallet.AccountNumber(sendPage.Sending.SelectedAccountLabel.Text)
	if err != nil {
		sendPage.showErrorLabel(values.AccountDetailsErr)
		return
	}

	balance, err := sendPage.Sending.SelectedWallet.GetAccountBalance(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		sendPage.showErrorLabel(values.AccountBalanceErr)
		return
	}

	sendPage.SpendableLabel.Text = values.SpendableAmountLabel + dcrutil.Amount(balance.Spendable).String()

	sendPage.amountEntry.OnChanged(sendPage.amountEntry.Text)
}

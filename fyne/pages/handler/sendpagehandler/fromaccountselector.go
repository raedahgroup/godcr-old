package sendpagehandler

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/handler/constantvalues"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initFromAccountSelector() error {
	fromLabel := widgets.NewTextWithStyle(constantvalues.FromText, color.RGBA{61, 88, 115, 255}, fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, constantvalues.DefaultTextSize)

	sendPage.Sending.OnAccountChange = sendPage.onAccountChange

	accountBox, err := sendPage.Sending.CreateAccountSelector(constantvalues.FromAccountSelectorPopUpHeaderLabel)
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
		sendPage.showErrorLabel(constantvalues.AccountDetailsErr)
		return
	}

	balance, err := sendPage.Sending.SelectedWallet.GetAccountBalance(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		sendPage.showErrorLabel(constantvalues.AccountBalanceErr)
		return
	}

	sendPage.SpendableLabel.Text = constantvalues.SpendableAmountLabel + dcrutil.Amount(balance.Spendable).String()
	sendPage.SpendableLabel.Refresh()

	sendPage.transactionFeeLabel.Refresh()
	sendPage.transactionSizeLabel.Refresh()
	sendPage.amountEntry.OnChanged(sendPage.amountEntry.Text)
}

package sendpagehandler

import (
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
)

func (sendPage *SendPageObjects) initFromAccountSelector() error {
	fromLabel := canvas.NewText(values.FromText, values.DarkerBlueGrayTextColor)

	sendPage.Sending.OnAccountChange = sendPage.onAccountChange

	sendPage.Sending.DefaultThemeColor = true
	accountBox, err := sendPage.Sending.CreateAccountSelector(values.FromAccountSelectorPopUpHeaderLabel)
	if err != nil { // return err if icons in account selector dont load
		return err
	}

	box := widget.NewVBox(
		fromLabel,
		accountBox)

	sendPage.SendPageContents.Append(box)
	return err
}

func (sendPage *SendPageObjects) onAccountChange() {
	balance, err := sendPage.Sending.SelectedWallet.GetAccountBalance(int32(*sendPage.Sending.SendingSelectedAccountID),
		dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		sendPage.showErrorLabel(values.AccountBalanceErr)
		return
	}

	sendPage.SendPageContents.Refresh()
	sendPage.SpendableLabel.Text = values.SpendableAmountLabel + dcrutil.Amount(balance.Spendable).String()
	sendPage.SendPageContents.Refresh()

	sendPage.SendPageContents.Refresh()
	if sendPage.sendMax {
		sendPage.maxButton.Container.OnTapped()
	} else {
		sendPage.initTxDetails(sendPage.amountEntry.Text)
	}
	sendPage.SendPageContents.Refresh()
}

package pages

import (
	"fmt"

	"fyne.io/fyne/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
)

type listener struct {
	tabMenu *widget.TabContainer
	dcrlw   *dcrlibwallet.LibWallet
	// Place widgets that should be dynamically updated on fyne.
	// Only update currently viewed tab and only update widget children.
}

func (test *listener) OnTransaction(transaction string) {
	if test.tabMenu.CurrentTabIndex() == 0 {
		// place overview page dynamic data here
	} else if test.tabMenu.CurrentTabIndex() == 2 { // place send page dynamic data here
		accountNumber, err := test.dcrlw.AccountNumber(sendPage.receivingSelectedAccountLabel.Text)
		if err != nil {
			return
		}
		balance, err := test.dcrlw.GetAccountBalance(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
		sendPage.receivingSelectedAccountBalanceLabel.SetText(dcrutil.Amount(balance.Total).String())

		accountNumber, err = test.dcrlw.AccountNumber(sendPage.sendingSelectedAccountLabel.Text)
		if err != nil {
			return
		}
		balance, err = test.dcrlw.GetAccountBalance(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
		sendPage.sendingSelectedAccountBalanceLabel.SetText(dcrutil.Amount(balance.Total).String())

		accounts, _ := test.dcrlw.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
		updateAccountDropdownContent(sendPage.sendingAccountDropdownContent, accounts)
	} else if test.tabMenu.CurrentTabIndex() == 3 {
		// place receive page dynamic data here
	} else if test.tabMenu.CurrentTabIndex() == 2 {
		// place account page dynamic data here
	}
}

func (test *listener) OnTransactionConfirmed(hash string, height int32) {
	fmt.Println("working OnTransactionConfirmed", hash, height)
	if test.tabMenu.CurrentTabIndex() == 0 {
		// place overview page dynamic data here
	} else if test.tabMenu.CurrentTabIndex() == 2 {
		// place send page dynamic data here

		accounts, _ := test.dcrlw.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
		updateAccountDropdownContent(sendPage.receivingAccountDropdownContent, accounts)
	} else if test.tabMenu.CurrentTabIndex() == 3 {
		// place receive page dynamic data here
	} else if test.tabMenu.CurrentTabIndex() == 2 {
		// place account page dynamic data here
	}
}

func (app *multiWalletTxListener) OnBlockAttached(walletID int, blockHeight int32) {
	// place all dynamic widgets in a function here, to be updated only when tabmenu is in view.
	if app.tabMenu.CurrentTabIndex() == 0 {

	} else if app.tabMenu.CurrentTabIndex() == 2 {

	} else if app.tabMenu.CurrentTabIndex() == 3 {

	} else if app.tabMenu.CurrentTabIndex() == 2 {

	}
}

func (app *AppInterface) walletNotificationListener(dcrlistener *listener) {
	// pass dynamic variables.
	dcrlistener.tabMenu = app.tabMenu
	dcrlistener.dcrlw = app.Dcrlw
	app.Dcrlw.TransactionNotification(dcrlistener)
}

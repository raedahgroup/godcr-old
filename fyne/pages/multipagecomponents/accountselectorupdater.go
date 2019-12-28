package multipagecomponents

import (
	"image/color"
	"log"
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func UpdateAccountSelectorOnNotification(accountBoxes []*widget.Box, sendingSelectedAccountBalanceLabel,
	spendableLabel *canvas.Text, multiWallet *dcrlibwallet.MultiWallet, selectedWalletID, selectedAccountID int, contents *widget.Box) {

	if contents == nil {
		return
	}

	selectedWalletIDs := multiWallet.OpenedWalletIDsRaw()
	sort.Ints(selectedWalletIDs)
	if len(selectedWalletIDs) != len(accountBoxes) {
		return
	}

	for walletIndex, accountBox := range accountBoxes {
		wallet := multiWallet.WalletWithID(selectedWalletIDs[walletIndex])
		if wallet == nil {
			return
		}

		account, err := wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
		if err != nil {
			log.Println("could not retrieve accounts on transaction notification")
			continue
		}

		if len(accountBox.Children) != len(account.Acc)-1 {
			continue
		}

		for index, boxContent := range accountBox.Children {
			spendableAmountLabel := widgets.NewTextWithStyle(dcrutil.Amount(account.Acc[index].Balance.Spendable).String(), color.Black, fyne.TextStyle{}, fyne.TextAlignTrailing, values.SpacerSize10)

			accountBalance := dcrutil.Amount(account.Acc[index].Balance.Total).String()
			accountBalanceLabel := widget.NewLabel(accountBalance)
			accountBalanceLabel.Alignment = fyne.TextAlignTrailing

			accountBalanceBox := widget.NewVBox(
				accountBalanceLabel,
				spendableAmountLabel,
			)

			if content, ok := boxContent.(*widgets.ClickableBox); ok {
				content.Box.Children[6] = accountBalanceBox
				content.Box.Refresh()
				content.Refresh()
			}
		}
	}

	wallet := multiWallet.WalletWithID(selectedWalletID)
	if wallet == nil {
		log.Println("could not retrieve selected wallet on transaction notification")
		return
	}

	account, err := wallet.GetAccount(int32(selectedAccountID), dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		log.Println("could not retrieve selected account on transaction notification")
		return
	}

	sendingSelectedAccountBalanceLabel.Text = dcrutil.Amount(account.TotalBalance).String()

	if spendableLabel != nil {
		spendableLabel.Text = values.SpendableAmountLabel + dcrutil.Amount(account.Balance.Spendable).String()
		canvas.Refresh(spendableLabel)
	}

	contents.Refresh()
}

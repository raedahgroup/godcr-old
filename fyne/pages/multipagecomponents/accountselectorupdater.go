package multipagecomponents

import (
	"image/color"
	"log"
	"sort"

	"github.com/raedahgroup/godcr/fyne/pages/constantvalues"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/widgets"
)

func UpdateAccountSelectorOnNotification(accountBoxes []*widget.Box, sendingSelectedAccountLabel, sendingSelectedAccountBalanceLabel *widget.Label,
	spendableLabel *canvas.Text, multiWallet *dcrlibwallet.MultiWallet, selectedWalletID int, contents *widget.Box) {

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
			spendableAmountLabel := canvas.NewText(dcrutil.Amount(account.Acc[index].Balance.Spendable).String(), color.Black)
			spendableAmountLabel.TextSize = 10
			spendableAmountLabel.Alignment = fyne.TextAlignTrailing

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

	accountNumber, err := wallet.AccountNumber(sendingSelectedAccountLabel.Text)
	if err != nil {
		log.Println("could not retrieve selected account number on transaction notification")
		return
	}

	account, err := wallet.GetAccount(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		log.Println("could not retrieve selected account on transaction notification")
		return
	}

	sendingSelectedAccountBalanceLabel.SetText(dcrutil.Amount(account.TotalBalance).String())

	if spendableLabel != nil {
		spendableLabel.Text = constantvalues.SpendableAmountLabel + dcrutil.Amount(account.Balance.Spendable).String()
		canvas.Refresh(spendableLabel)
	}

	contents.Refresh()
}

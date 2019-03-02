package pages

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func SendPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)
	// textView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	form := tview.NewForm()

	body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Send Fund"), 4, 1, false)
	
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}

	//Form for Sending
	var  amount, destination []string
	var Passphrase string
	var accountNum uint32
	var spendUnconfirmed bool

	accountNames := make([]string, len(accounts))
	accountN := make([]uint32, len(accounts))
	for index, account := range accounts {
		accountNames[index] = fmt.Sprintf("%s - %s ", account.Name, account.Balance.Total.String())
		body.AddItem(form.AddDropDown("Account", []string{accountNames[index]}, 0, func(option string, optionIndex int) {
			accountNum = accountN[optionIndex]
		}).
		AddInputField("Amount", "", 20, nil, func (text string) {
			amount = text
		}).
		AddInputField("Destination Address", "", 20, nil, func (text string) {
			destination = text
		}).
		AddCheckbox("Unconfirmed", false, func(checked bool) {
			if checked {
				spendUnconfirmed = true
			}
		}).
		AddCheckbox("Select custom inputs", false, func(checked bool) {
			//todo add select custom inputs to send page
		}).
		AddPasswordField("Wallet Passphrase", "", 20, '*', func (text string) {
			Passphrase = text
		}).
		AddButton("Send", func() {
			sendDestinations := make([]txhelper.TransactionDestination, len(destination))
			for i := range destination {
				Amount, err := strconv.ParseFloat(string(amount[i]), 32)
				if err != nil {
					body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
					return
				}
				sendDestinations[i] = txhelper.TransactionDestination{
					Address: destination[i],
					Amount:  amount,
				}
			}

			var requiredConfirmations int32 = walletcore.DefaultRequiredConfirmations
			if spendUnconfirmed {
				requiredConfirmations = 0
			}

			txHash, err = routes.walletMiddleware.SendFromAccount(accountNum, requiredConfirmations, sendDestinations, passphrase)
			if err != nil {
				body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
				return
			}
			fmt.Println(txHash)
		}).
		AddButton("Cancel", func() {
			clearFocus()
		}), 0, 1, true)
	}

	setFocus(body)
	
	return body
}

// func SendPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
// 	textView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
// 	flex := tview.NewFlex()
// }

func confBalance(accountNum uint32, wallet walletcore.Wallet) error{
	accountBalance, err := wallet.AccountBalance(accountNum, walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return err
	}
	if accountBalance.Total != 0 {
		return fmt.Errorf("Selected account has 0 balance. Cannot proceed")
	}

	fmt.Println(accountBalance)
	return nil
}
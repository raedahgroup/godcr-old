package pages

import (
	"fmt"
	"strconv"

	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/rivo/tview"
)

func SendPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Send Fund"), 4, 1, false)

	textView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		body.RemoveItem(textView)
		return body.AddItem(textView.SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}

	//Form for Sending
	form := tview.NewForm()

	accountNames := make([]string, len(accounts))
	accountNumber := make([]uint32, len(accounts))
	for index, account := range accounts {
		accountNames[index] = fmt.Sprintf("%s - %s ", account.Name, account.Balance.Total.String())
		accountNumber[index] = account.Number
	}

	var accountNum uint32
	form.AddDropDown("Account", accountNames, 0, func(option string, optionIndex int) {
		accountNum = accountNumber[optionIndex]
	})
	var amount string
	form.AddInputField("Amount", "", 20, nil, func(text string) {
		amount = text
	})

	var destination string
	form.AddInputField("Destination Address", "", 30, nil, func(text string) {
		destination = text
	})

	var spendUnconfirmed bool
	form.AddCheckbox("spend Unconfirmed", false, func(checked bool) {
		if checked {
			spendUnconfirmed = true
		}
	})

	form.AddCheckbox("Select custom inputs", false, func(checked bool) {
		//todo add select custom inputs feature to send page
	})

	var passphrase string
	form.AddPasswordField("Wallet Passphrase", "", 20, '*', func(text string) {
		passphrase = text
	})
	form.AddButton("Send", func() {
		sendDestinations := make([]txhelper.TransactionDestination, len(destination))
		for i := range destination {
			Amount, err := strconv.ParseFloat(string(amount), 64)
			if err != nil {
				body.RemoveItem(textView)
				body.AddItem(textView.SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
				return
			}
			sendDestinations[i] = txhelper.TransactionDestination{
				Address: destination,
				Amount:  Amount,
			}
		}

		var requiredConfirmations int32 = walletcore.DefaultRequiredConfirmations
		if spendUnconfirmed {
			requiredConfirmations = 0
		}

		txHash, err := wallet.SendFromAccount(accountNum, requiredConfirmations, sendDestinations, passphrase)
		if err != nil {
			body.RemoveItem(textView)
			body.AddItem(textView.SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
			return
		}
		body.RemoveItem(textView)
		body.AddItem(textView.SetText(fmt.Sprintf("Sent txid", txHash)), 0, 1, false)

	})
	form.AddButton("Cancel", func() {
		clearFocus()
	})

	body.AddItem(form, 16, 1, true)

	setFocus(body)

	return body
}

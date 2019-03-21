package pages

import (
	"fmt"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"strconv"

	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/rivo/tview"
)

func SendPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Send"), 2, 1, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}

	accountNames := make([]string, len(accounts))
	accountNumbers := make([]uint32, len(accounts))
	for index, account := range accounts {
		accountNames[index] = account.String()
		accountNumbers[index] = account.Number
	}

	// Form for Sending
	form := primitives.NewForm()
	var accountNum uint32
	form.AddDropDown("Source Account", accountNames, 0, func(option string, optionIndex int) {
		accountNum = accountNumbers[optionIndex]
	})

	var amount string
	form.AddInputField("Amount", "", 20, nil, func(text string) {
		amount = text
	})

	var destination string
	form.AddInputField("Destination Address", "", 40, nil, func(text string) {
		destination = text
	})

	var spendUnconfirmed bool
	form.AddCheckbox("Spend Unconfirmed", false, func(checked bool) {
		spendUnconfirmed = checked
	})

	var passphrase string
	form.AddPasswordField("Wallet Passphrase", "", 20, '*', func(text string) {
		passphrase = text
	})

	outputTextView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	outputMessage := func(output string) {
		body.RemoveItem(outputTextView)
		body.AddItem(outputTextView.SetText(output), 0, 1, true)
	}

	form.AddButton("Send", func() {
		sendDestination := make([]txhelper.TransactionDestination, 1)
		amount, err := strconv.ParseFloat(string(amount), 64)
		if err != nil {
			outputMessage("Error: Invalid amount")
			return
		}
		sendDestination[0] = txhelper.TransactionDestination{
			Address: destination,
			Amount:  amount,
		}

		var requiredConfirmations int32 = walletcore.DefaultRequiredConfirmations
		if spendUnconfirmed {
			requiredConfirmations = 0
		}

		txHash, err := wallet.SendFromAccount(accountNum, requiredConfirmations, sendDestination, passphrase)
		if err != nil {
			outputMessage(err.Error())
			return
		}

		outputMessage("Sent txid " + txHash)
	})

	form.SetCancelFunc(clearFocus)
	body.AddItem(form, 13, 1, true)

	setFocus(body)

	return body
}

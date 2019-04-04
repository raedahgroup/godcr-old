package pages

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func sendPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)
	
	hintTextView.SetText("TIP: Select source account with ARROW DOWN and ENTER. Move around with TAB and SHIFT+TAB. ESC to return to Navigation menu")

	// page title and tip
	body.AddItem(primitives.NewLeftAlignedTextView("SEND"), 2, 0, false)
	
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}

	// form for Sending
	form := primitives.NewForm()
	form.SetBorderPadding(0, 0, 0, 0)
	body.AddItem(form, 0, 1, true)

	outputMessageTextView := primitives.WordWrappedTextView("")
	outputMessageTextView.SetTextColor(tcell.ColorOrangeRed)

	displayMessage := func(message string) {
		body.RemoveItem(outputMessageTextView)
		outputMessageTextView.SetText(message)
		body.AddItem(outputMessageTextView, 2, 0, false)
	}
	
	accountNames := make([]string, len(accounts))
	accountNumbers := make([]uint32, len(accounts))
	for index, account := range accounts {
		accountNames[index] = account.String()
		accountNumbers[index] = account.Number
	}

	// add form fields
	var accountNum uint32
	form.AddDropDown("Source Account", accountNames, 0, func(option string, optionIndex int) {
		accountNum = accountNumbers[optionIndex]
	})

	var amount string
	form.AddInputField("Amount", "", 20, nil, func(text string) {
		amount = text
	})

	var destination string
	form.AddInputField("Destination Address", "", 37, nil, func(text string) {
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

	form.AddButton("Send", func() {
		amount, err := strconv.ParseFloat(string(amount), 64)
		if err != nil {
			displayMessage("Error: Invalid amount")
			return
		}

		sendDestination := make([]txhelper.TransactionDestination, 1)
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
			displayMessage(err.Error())
			return
		}

		displayMessage("Sent txid " + txHash)

		// reset form
		form.ClearFields()
		setFocus(form.GetFormItem(0))
	})

	form.SetCancelFunc(clearFocus)

	setFocus(body)

	return body
}

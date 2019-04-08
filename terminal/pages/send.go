package pages

import (
	"fmt"
	"strconv"

	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func sendPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	pages := tview.NewPages()

	body := tview.NewFlex().SetDirection(tview.FlexRow)
	pages.AddPage("main", body, true, true)

	body.AddItem(primitives.NewLeftAlignedTextView("Sending Decred"), 2, 0, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}

	// form for Sending
	form := primitives.NewForm()
	form.SetBorderPadding(0, 0, 0, 0)
	body.AddItem(form, 0, 1, true)

	errorTextView := primitives.WordWrappedTextView("")
	errorTextView.SetTextColor(helpers.DecredOrangeColor)

	displayErrorMessage := func(message string) {
		body.RemoveItem(errorTextView)
		errorTextView.SetText(message)
		body.AddItem(errorTextView, 2, 0, false)
	}

	accountNames := make([]string, len(accounts))
	accountNumbers := make([]uint32, len(accounts))
	for index, account := range accounts {
		accountNames[index] = account.String()
		accountNumbers[index] = account.Number
	}

	// add form fields
	var accountNum uint32
	form.AddDropDown("From", accountNames, 0, func(option string, optionIndex int) {
		accountNum = accountNumbers[optionIndex]
	})

	var destination string
	form.AddInputField("Destination Address", "", 37, nil, func(text string) {
		destination = text
	})

	var amount string
	form.AddInputField("Amount", "", 20, nil, func(text string) {
		amount = text
	})
	
	var spendUnconfirmed bool
	form.AddCheckbox("Spend Unconfirmed", false, func(checked bool) {
		spendUnconfirmed = checked
	})

	form.AddButton("Send", func() {
		// validate form fields
		amount, err := strconv.ParseFloat(string(amount), 64)
		if err != nil {
			displayErrorMessage("Error: Invalid amount")
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

		helpers.RequestSpendingPassphrase(pages, func(passphrase string) {
			setFocus(form)

			txHash, err := wallet.SendFromAccount(accountNum, requiredConfirmations, sendDestination, passphrase)
			if err != nil {
				displayErrorMessage(err.Error())
				return
			}

			body.AddItem(primitives.WordWrappedTextView("Sent txid "+txHash), 2, 0, false)

			// reset form
			form.ClearFields()
			setFocus(form.GetFormItem(0))
		}, func() {
			setFocus(form)
		})
	})

	form.SetCancelFunc(clearFocus)

	hintTextView.SetText("TIP: Select source account with ARROW DOWN and ENTER. Move around with TAB and SHIFT+TAB. ESC to return to navigation menu")

	setFocus(pages)

	return pages
}

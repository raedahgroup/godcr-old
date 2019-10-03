package pages

import (
	"fmt"
	"strconv"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func sendPage() tview.Primitive {
	pages := tview.NewPages()

	body := tview.NewFlex().SetDirection(tview.FlexRow)
	pages.AddPage("main", body, true, true)

	body.AddItem(primitives.NewLeftAlignedTextView("Sending Decred"), 2, 0, false)

	getAccountsResp, err := commonPageData.wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}

	// form for Sending
	form := primitives.NewForm(true)
	form.SetBorderPadding(0, 0, 0, 0)
	body.AddItem(form, 0, 1, true)

	errorTextView := primitives.WordWrappedTextView("")
	errorTextView.SetTextColor(helpers.DecredOrangeColor)

	displayErrorMessage := func(message string) {
		body.RemoveItem(errorTextView)
		errorTextView.SetText(message)
		body.AddItem(errorTextView, 2, 0, false)
	}

	accountSelectionWidgetData := &helpers.AccountSelectionWidgetData{
		Label:    "From:",
		Accounts: getAccountsResp.Acc,
	}
	helpers.AddAccountSelectionWidgetToForm(form, accountSelectionWidgetData)

	var destination string
	form.AddInputField("Destination Address:", "", 37, nil, func(text string) {
		destination = text
	})

	var amount string
	form.AddInputField("Amount:", "", 20, nil, func(text string) {
		amount = text
	})

	var spendUnconfirmed bool
	form.AddCheckbox("Spend Unconfirmed:", false, func(checked bool) {
		spendUnconfirmed = checked
	})

	form.AddButton("Send", func() {
		// validate form fields
		amount, err := strconv.ParseFloat(string(amount), 64)
		if err != nil {
			displayErrorMessage("Error: Invalid amount")
			return
		}

		var requiredConfirmations int32 = dcrlibwallet.DefaultRequiredConfirmations
		if spendUnconfirmed {
			requiredConfirmations = 0
		}

		helpers.RequestSpendingPassphrase(pages, func(passphrase string) {
			// return focus to the form
			commonPageData.app.SetFocus(form)

			// create and broadcast new tx
			newTx := commonPageData.wallet.NewUnsignedTx(accountSelectionWidgetData.SelectedAccountNumber, requiredConfirmations)
			newTx.AddSendDestination(destination, dcrlibwallet.AmountAtom(amount), false)

			txHash, err := newTx.Broadcast([]byte(passphrase))
			if err != nil {
				displayErrorMessage(err.Error())
				return
			}

			body.AddItem(primitives.WordWrappedTextView("Sent txid "+dcrlibwallet.EncodeHex(txHash)), 2, 0, false)

			// reset form
			form.ClearFields()
			commonPageData.app.SetFocus(form.GetFormItem(0))
		}, func() {
			commonPageData.app.SetFocus(form)
		})
	})

	form.AddButton("Clear", func() {
		form.ClearFields()
		body.RemoveItem(errorTextView)
	})

	form.SetCancelFunc(commonPageData.clearAllPageContent)

	if len(getAccountsResp.Acc) <= 1 {
		commonPageData.hintTextView.SetText("TIP: Move around with TAB and SHIFT+TAB. ESC to return to navigation menu")
	} else {
		commonPageData.hintTextView.SetText("TIP: Select source account with ARROW DOWN and ENTER. Move around with TAB. ESC to return to navigation menu")
	}

	commonPageData.app.SetFocus(pages)

	return pages
}

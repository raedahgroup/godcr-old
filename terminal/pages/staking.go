package pages

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func stakingPage() tview.Primitive {
	// parent flexbox layout container to hold other primitives
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	body.AddItem(primitives.NewLeftAlignedTextView("Staking"), 2, 0, false)

	messageTextView := primitives.WordWrappedTextView("")

	clearMessage := func() {
		body.RemoveItem(messageTextView)
	}

	displayMessage := func(message string, error bool) {
		clearMessage()
		messageTextView.SetText(message)
		if error {
			messageTextView.SetTextColor(helpers.DecredOrangeColor)
		} else {
			messageTextView.SetTextColor(helpers.DecredGreenColor)
		}
		body.AddItem(messageTextView, 2, 0, false)
	}

	stakeInfo, err := stakeInfoFlex()
	if err != nil {
		errorText := fmt.Sprintf("Error fetching stake info: %s", err.Error())
		displayMessage(errorText, true)
	} else {
		body.AddItem(stakeInfo, 3, 0, false)
	}

	body.AddItem(tview.NewTextView().SetText("-Purchase Ticket-").SetTextColor(helpers.DecredLightBlueColor), 2, 0, false)
	purchaseTicket, err := purchaseTicketForm(displayMessage, clearMessage)
	if err != nil {
		errorText := fmt.Sprintf("Error setting up purchase form: %s", err.Error())
		displayMessage(errorText, true)
	} else {
		body.AddItem(purchaseTicket, 0, 1, true)
	}

	commonPageData.app.SetFocus(body)

	commonPageData.hintTextView.SetText("TIP: Move around with TAB and SHIFT+TAB. ESC to return to navigation menu")

	return body
}

func stakeInfoFlex() (*primitives.TextView, error) {
	stakeInfo, err := commonPageData.wallet.StakeInfo()
	if err != nil {
		return nil, err
	} else if stakeInfo == nil {
		return nil, errors.New("no tickets in wallet")
	}

	stakingReport := fmt.Sprintf("Mempool: %d  Immature: %d  Live: %d", stakeInfo.OwnMempoolTix, stakeInfo.Immature, stakeInfo.Live)
	return primitives.NewLeftAlignedTextView(stakingReport), nil
}

func purchaseTicketForm(displayMessage func(message string, error bool), clearMessage func()) (*tview.Pages, error) {
	pages := tview.NewPages()

	getAccountsResp, err := commonPageData.wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return nil, err
	}

	form := primitives.NewForm(true)
	form.SetBorderPadding(0, 0, 0, 0)
	pages.AddPage("form", form, true, true)

	accountSelectionWidgetData := &helpers.AccountSelectionWidgetData{
		Label:    "From:",
		Accounts: getAccountsResp.Acc,
	}
	helpers.AddAccountSelectionWidgetToForm(form, accountSelectionWidgetData)

	var numTickets string
	form.AddInputField("Number of Tickets:", "", 10, nil, func(text string) {
		numTickets = text
	})

	var spendUnconfirmed bool
	form.AddCheckbox("Spend Unconfirmed:", false, func(checked bool) {
		spendUnconfirmed = checked
	})

	form.AddButton("Purchase", func() {
		if len(numTickets) == 0 {
			displayMessage("Error: please specify the number of tickets to purchase", true)
			return
		}

		helpers.RequestSpendingPassphrase(pages, func(passphrase string) {
			commonPageData.app.SetFocus(form)

			accountNumber := accountSelectionWidgetData.SelectedAccountNumber
			ticketHashes, err := purchaseTickets(passphrase, numTickets, accountNumber, spendUnconfirmed)
			if err != nil {
				displayMessage(err.Error(), true)
				return
			}

			successMessage := fmt.Sprintf("You have purchased %d ticket(s)\n%s", len(ticketHashes), strings.Join(ticketHashes, "\n"))
			displayMessage(successMessage, false)

			// reset form
			form.ClearFields()
			commonPageData.app.SetFocus(form.GetFormItem(0))
		}, func() {
			commonPageData.app.SetFocus(form)
		})
	})

	form.AddButton("Clear", func() {
		form.ClearFields()
		clearMessage()
	})

	form.SetCancelFunc(commonPageData.clearAllPageContent)

	return pages, nil
}

func purchaseTickets(passphrase, numTickets string, accountNum int32, spendUnconfirmed bool) ([]string, error) {
	nTickets, err := strconv.ParseUint(string(numTickets), 10, 32)
	if err != nil {
		return nil, err
	}

	requiredConfirmations := dcrlibwallet.DefaultRequiredConfirmations
	if spendUnconfirmed {
		requiredConfirmations = 0
	}

	request := &dcrlibwallet.PurchaseTicketsRequest{
		RequiredConfirmations: uint32(requiredConfirmations),
		Passphrase:            []byte(passphrase),
		NumTickets:            uint32(nTickets),
		Account:               uint32(accountNum),
	}

	ticketHashes, err := commonPageData.wallet.PurchaseTickets(context.Background(), request)
	if err != nil {
		return nil, err
	}

	if len(ticketHashes) == 0 {
		err = errors.New("no ticket was purchased")
		return nil, err
	}

	return ticketHashes, nil
}

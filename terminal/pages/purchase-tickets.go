package pages

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/rivo/tview"
)

func PurchaseTicketsPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Purchase Tickets"), 4, 1, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}

	accountNumbers := make([]uint32, len(accounts))
	accountOverviews := make([]string, len(accounts))
	for index, account := range accounts {
		accountOverviews[index] = fmt.Sprintf("%s - Total %s (Spendable %s)", account.Name, account.Balance.Total.String(), account.Balance.Spendable.String())
		accountNumbers[index] = account.Number
	}

	form := tview.NewForm()
	var accountNum uint32
	form.AddDropDown("Source Account", accountOverviews, 0, func(option string, optionIndex int) {
		accountNum = accountNumbers[optionIndex]
	})

	var numTickets string
	form.AddInputField("Number of tickets", "", 20, nil, func(text string) {
		numTickets = text
	})

	var spendUnconfirmed bool
	form.AddCheckbox("Spend Unconfirmed", false, func(checked bool) {
		spendUnconfirmed = checked
	})

	var passphrase string
	form.AddPasswordField("Spending Passphrase", "", 20, '*', func(text string) {
		passphrase = text
	})

	outputTextView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	outputMessage := func(ticketHashes []string, err error) {
		body.RemoveItem(outputTextView)
		if err != nil {
			body.AddItem(outputTextView.SetText(fmt.Sprintf(err.Error())), 0, 1, true)
		} else {
			body.AddItem(outputTextView.SetText(fmt.Sprintf("You have purchased %d ticket(s)\n%s", len(ticketHashes), strings.Join(ticketHashes, "\n"))), 0, 1, true)
		}
	}

	form.AddButton("Submit", func() {
		if len(numTickets) == 0 {
			err := errors.New("Error: please specify the number of tickets to purchase")
			outputMessage(nil, err)
			return
		}
		if len(passphrase) == 0 {
			err := errors.New("Error: please enter your spending passphrase")
			outputMessage(nil, err)
			return
		}

		ticketHashes, err := submit(passphrase, numTickets, accountNum, spendUnconfirmed, wallet)
		if err != nil {
			outputMessage(nil, err)
			return
		}

		outputMessage(ticketHashes, nil)
	})

	form.SetCancelFunc(clearFocus)

	body.AddItem(form, 12, 1, true)

	setFocus(body)
	return body
}

func submit(passphrase, numTickets string, accountNum uint32, spendUnconfirmed bool, wallet walletcore.Wallet) ([]string, error) {
	nTickets, err := strconv.ParseUint(string(numTickets), 10, 32)
	if err != nil {
		return nil, err
	}

	requiredConfirmations := walletcore.DefaultRequiredConfirmations
	if spendUnconfirmed {
		requiredConfirmations = 0
	}

	request := dcrlibwallet.PurchaseTicketsRequest{
		RequiredConfirmations: uint32(requiredConfirmations),
		Passphrase:            []byte(passphrase),
		NumTickets:            uint32(nTickets),
		Account:               accountNum,
	}

	ticketHashes, err := wallet.PurchaseTickets(context.Background(), request)
	if err != nil {
		return nil, err
	}

	if len(ticketHashes) == 0 {
		err = errors.New("no ticket was purchased")
		return nil, err
	}

	return ticketHashes, nil
}

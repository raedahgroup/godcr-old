package pages

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/rivo/tview"
)

func PurchaseTicketPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)
	textView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Purchase Tickets"), 4, 1, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return body.AddItem(textView.SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}

	accountNumbers := make([]uint32, len(accounts))
	accountOverviews := make([]string, len(accounts))
	for index, account := range accounts {
		accountOverviews[index] = fmt.Sprintf("%s - Total %s (Spendable %s)", account.Name, account.Balance.Total.String(), account.Balance.Spendable.String())
		accountNumbers[index] = account.Number
	}

	var accountNum uint32
	var numTickets, passphrase string
	var spendUnconfirmed bool
	form := tview.NewForm()
	body.AddItem(form.AddDropDown("Source Account", accountOverviews, 0, func(option string, optionIndex int) {
		accountNum = accountNumbers[optionIndex]
	}).
		AddInputField("Number of tickets", "", 20, nil, func(text string) {
			numTickets = text
		}).
		AddCheckbox("Spend Unconfirmed", false, func(checked bool) {
			if checked {
				spendUnconfirmed = true
			}
		}).
		AddPasswordField("wallet Passphrase", "", 20, '*', func(text string) {
			passphrase = text
		}).
		AddButton("Submit", func() {
			if len(numTickets) == 0 {
				body.RemoveItem(textView)
				body.AddItem(textView.SetText(fmt.Sprintf("Error: %s", "Please specify the number of tickets to purchase")).SetDoneFunc(func(key tcell.Key) {
				}), 0, 1, true)
				return
			}
			if len(passphrase) == 0 {
				body.RemoveItem(textView)
				body.AddItem(textView.SetText(fmt.Sprintf("Error: %s", "please enter your wallet passphrase")).SetDoneFunc(func(key tcell.Key) {
					if key == tcell.KeyEscape {
						clearFocus()
					}
				}), 0, 1, true)
				return
			}

			ticketHashes, err := submit(passphrase, numTickets, accountNum, spendUnconfirmed, wallet)
			if err != nil {
				body.RemoveItem(textView)
				body.AddItem(textView.SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, true)
				return
			}
			body.RemoveItem(textView)
			body.AddItem(textView.SetText(fmt.Sprintf("You have purchased %d ticket(s)\n%s", len(ticketHashes), strings.Join(ticketHashes, "\n"))), 0, 1, true)
		}).SetCancelFunc(func() {
		clearFocus()
	}), 12, 1, true)

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
		Account:               uint32(accountNum),
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

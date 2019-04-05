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
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func stakingPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	// parent flexbox layout container to hold other primitives
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	// page title and tip
	body.AddItem(primitives.NewLeftAlignedTextView("STAKING"), 2, 0, false)
	
	errorTextView := primitives.WordWrappedTextView("")
	errorTextView.SetTextColor(tcell.ColorOrangeRed)

	displayError := func(errorMessage string) {
		body.RemoveItem(errorTextView)
		errorTextView.SetText(errorMessage)
		body.AddItem(errorTextView, 0, 1, false)
	}

	body.AddItem(tview.NewTextView().SetText("Stake Info").SetTextColor(helpers.DecredLightColor), 1, 0, false)
	stakeInfo, err := stakeInfoTable(wallet)
	if err != nil {
		errorText := fmt.Sprintf("Error fetching stake info: %s", err.Error())
		displayError(errorText)
	} else {
		body.AddItem(stakeInfo, 6, 0, true)
	}

	body.AddItem(tview.NewTextView().SetText("Purchase Ticket").SetTextColor(helpers.DecredLightColor), 2, 0, false)
	purchaseTicket, statusTextView, err := purchaseTicketForm(wallet, displayError)
	if err != nil {
		errorText := fmt.Sprintf("Error setting up purchase form: %s", err.Error())
		displayError(errorText)
	} else {
		body.AddItem(purchaseTicket, 12, 0, true)
		body.AddItem(statusTextView, 3, 0, true)
	}

	stakeInfo.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyTAB {
			setFocus(purchaseTicket)
			setFocus(purchaseTicket.GetFormItem(0))
			return nil
		}

		return event
	})

	// listen to escape and left key press events on all form items and buttons
	purchaseTicket.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			clearFocus()
			return nil
		}
		return event
	})

	// use different key press listener on first form item to watch for backtab press and restore focus to stake info
	purchaseTicket.GetFormItemBox(0).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			setFocus(stakeInfo)
			return nil
		}
		return event
	})

	// use different key press listener on form button to watch for tab press and restore focus to stake info
	purchaseTicket.GetButton(0).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyTAB {
			setFocus(stakeInfo)
			return nil
		}
		return event
	})

	setFocus(body)
	
	hintTextView.SetText("TIP: Move around with TAB and SHIFT+TAB. ESC to return to Navigation menu")

	return body
}

func stakeInfoTable(wallet walletcore.Wallet) (*primitives.Table, error) {
	stakeInfo, err := wallet.StakeInfo(context.Background())
	if err != nil {
		return nil, err
	} else if stakeInfo == nil {
		return nil, errors.New("no tickets in wallet")
	}

	table := primitives.NewTable()
	table.SetBorders(true)

	table.SetHeaderCell(0, 0, "Expired")
	table.SetHeaderCell(0, 1, "Immature")
	table.SetHeaderCell(0, 2, "Live")
	table.SetHeaderCell(0, 3, "Revoked")
	table.SetHeaderCell(0, 4, "Unmined")
	table.SetHeaderCell(0, 5, "Unspent")
	table.SetHeaderCell(0, 6, "AllmempoolTix")
	table.SetHeaderCell(0, 7, "PoolSize")
	table.SetHeaderCell(0, 8, "Missed")
	table.SetHeaderCell(0, 9, "Voted")
	table.SetHeaderCell(0, 10, "Total Subsidy")

	numberToString := func(n uint32) string {
		return strconv.Itoa(int(n))
	}

	table.SetCellCenterAlign(1, 0, numberToString(stakeInfo.Expired))
	table.SetCellCenterAlign(1, 1, numberToString(stakeInfo.Immature))
	table.SetCellCenterAlign(1, 2, numberToString(stakeInfo.Live))
	table.SetCellCenterAlign(1, 3, numberToString(stakeInfo.Revoked))
	table.SetCellCenterAlign(1, 4, numberToString(stakeInfo.OwnMempoolTix))
	table.SetCellCenterAlign(1, 5, numberToString(stakeInfo.Unspent))
	table.SetCellCenterAlign(1, 6, numberToString(stakeInfo.AllMempoolTix))
	table.SetCellCenterAlign(1, 7, numberToString(stakeInfo.PoolSize))
	table.SetCellCenterAlign(1, 8, numberToString(stakeInfo.Missed))
	table.SetCellCenterAlign(1, 9, numberToString(stakeInfo.Voted))
	table.SetCellCenterAlign(1, 10, stakeInfo.TotalSubsidy)

	return table, nil
}

func purchaseTicketForm(wallet walletcore.Wallet, displayError func(errorMessage string)) (*primitives.Form, *primitives.TextView, error) {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return nil, nil, err
	}

	accountNumbers := make([]uint32, len(accounts))
	accountOverviews := make([]string, len(accounts))
	for index, account := range accounts {
		accountOverviews[index] = account.String()
		accountNumbers[index] = account.Number
	}

	form := primitives.NewForm()
	form.SetBorderPadding(0, 0, 0, 0)

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

	// empty status text view for updating status of ticket purchase operation
	statusTextView := primitives.WordWrappedTextView("")

	form.AddButton("Submit", func() {
		if len(numTickets) == 0 {
			displayError("Error: please specify the number of tickets to purchase")
			return
		}
		if len(passphrase) == 0 {
			displayError("Error: please enter your spending passphrase")
			return
		}

		ticketHashes, err := purchaseTickets(passphrase, numTickets, accountNum, spendUnconfirmed, wallet)
		if err != nil {
			displayError(err.Error())
			return
		}

		output := fmt.Sprintf("You have purchased %d ticket(s)\n%s", len(ticketHashes), strings.Join(ticketHashes, "\n"))
		statusTextView.SetText(output)
	})

	return form, statusTextView, nil
}

func purchaseTickets(passphrase, numTickets string, accountNum uint32, spendUnconfirmed bool, wallet walletcore.Wallet) ([]string, error) {
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

	ticketHashes, err := wallet.PurchaseTicket(context.Background(), request)
	if err != nil {
		return nil, err
	}

	if len(ticketHashes) == 0 {
		err = errors.New("no ticket was purchased")
		return nil, err
	}

	return ticketHashes, nil
}

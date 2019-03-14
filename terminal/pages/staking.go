package pages

import (
	"context"
	"errors"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
	"strconv"
	"strings"
)

func StakingPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	// parent flexbox layout container to hold other primitives
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	// page title and tip
	body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Staking"), 1, 0, false)
	hintText := primitives.WordWrappedTextView("(TIP: Move around with Tab and Shift+Tab. Return to nav menu with Esc)")
	hintText.SetTextColor(tcell.ColorGray)
	body.AddItem(hintText, 2, 0, false)

	body.AddItem(tview.NewTextView().SetText("Stake Info").SetTextColor(helpers.DecredLightColor), 1, 0, false)
	stakeInfo, err := stakeInfoTable(wallet)
	if err != nil {
		errorText := fmt.Sprintf("Error fetching stake info: %s", err.Error())
		body.AddItem(primitives.WordWrappedTextView(errorText), 1, 0, false)
	} else {
		body.AddItem(stakeInfo, 6, 0, true)
	}

	body.AddItem(tview.NewTextView().SetText("Purchase Ticket").SetTextColor(helpers.DecredLightColor), 1, 0, false)
	purchaseTicket, err := purchaseTicketForm(wallet, body)
	if err != nil {
		errorText := fmt.Sprintf("Error setting up purchase form: %s", err.Error())
		body.AddItem(primitives.WordWrappedTextView(errorText), 1, 0, false)
	} else {
		body.AddItem(purchaseTicket, 0, 1, true)
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

	return body
}

func stakeInfoTable(wallet walletcore.Wallet) (*tview.Table, error) {
	stakeInfo, err := wallet.StakeInfo(context.Background())
	if err != nil {
		return nil, err
	} else if stakeInfo == nil {
		return nil, errors.New("no tickets in wallet")
	}

	table := tview.NewTable().SetBorders(true)

	table.SetCell(0, 0, tview.NewTableCell("Expired").SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("Immature").SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("Live").SetAlign(tview.AlignCenter))
	table.SetCell(0, 3, tview.NewTableCell("Revoked").SetAlign(tview.AlignCenter))
	table.SetCell(0, 4, tview.NewTableCell("Unmined").SetAlign(tview.AlignCenter))
	table.SetCell(0, 5, tview.NewTableCell("Unspent").SetAlign(tview.AlignCenter))
	table.SetCell(0, 6, tview.NewTableCell("AllmempoolTix").SetAlign(tview.AlignCenter))
	table.SetCell(0, 7, tview.NewTableCell("PoolSize").SetAlign(tview.AlignCenter))
	table.SetCell(0, 8, tview.NewTableCell("Missed").SetAlign(tview.AlignCenter))
	table.SetCell(0, 9, tview.NewTableCell("Voted").SetAlign(tview.AlignCenter))
	table.SetCell(0, 10, tview.NewTableCell("Total Subsidy").SetAlign(tview.AlignCenter))

	table.SetCell(1, 0, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Expired))).SetAlign(tview.AlignCenter))
	table.SetCell(1, 1, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Immature))).SetAlign(tview.AlignCenter))
	table.SetCell(1, 2, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Live))).SetAlign(tview.AlignCenter))
	table.SetCell(1, 3, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Revoked))).SetAlign(tview.AlignCenter))
	table.SetCell(1, 4, tview.NewTableCell(strconv.Itoa(int(stakeInfo.OwnMempoolTix))).SetAlign(tview.AlignCenter))
	table.SetCell(1, 5, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Unspent))).SetAlign(tview.AlignCenter))
	table.SetCell(1, 6, tview.NewTableCell(strconv.Itoa(int(stakeInfo.AllMempoolTix))).SetAlign(tview.AlignCenter))
	table.SetCell(1, 7, tview.NewTableCell(strconv.Itoa(int(stakeInfo.PoolSize))).SetAlign(tview.AlignCenter))
	table.SetCell(1, 8, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Missed))).SetAlign(tview.AlignCenter))
	table.SetCell(1, 9, tview.NewTableCell(strconv.Itoa(int(stakeInfo.Voted))).SetAlign(tview.AlignCenter))
	table.SetCell(1, 10, tview.NewTableCell(stakeInfo.TotalSubsidy).SetAlign(tview.AlignCenter))

	return table, nil
}

func purchaseTicketForm(wallet walletcore.Wallet, parentBody *tview.Flex) (*primitives.Form, error) {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return nil, err
	}

	accountNumbers := make([]uint32, len(accounts))
	accountOverviews := make([]string, len(accounts))
	for index, account := range accounts {
		accountOverviews[index] = account.String()
		accountNumbers[index] = account.Number
	}

	form := primitives.NewForm()
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
	outputMessage := func(output string) {
		parentBody.RemoveItem(outputTextView)
		parentBody.AddItem(outputTextView.SetText(output), 0, 1, true)
	}

	form.AddButton("Submit", func() {
		if len(numTickets) == 0 {
			outputMessage("Error: please specify the number of tickets to purchase")
			return
		}
		if len(passphrase) == 0 {
			outputMessage("Error: please enter your spending passphrase")
			return
		}

		ticketHashes, err := purchaseTickets(passphrase, numTickets, accountNum, spendUnconfirmed, wallet)
		if err != nil {
			outputMessage(err.Error())
			return
		}

		output := fmt.Sprintf("You have purchased %d ticket(s)\n%s", len(ticketHashes), strings.Join(ticketHashes, "\n"))
		outputMessage(output)
	})

	return form, nil
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

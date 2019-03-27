package pages

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func balancePage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	textView := tview.NewTextView()

	body := tview.NewFlex().SetDirection(tview.FlexRow)
	body.SetBorderPadding(1, 0, 1, 0)

	hintText := primitives.WordWrappedTextView("(TIP: Hit ENTER to switch between Detailed Balance and Simple Balance, ARROW keys to Scroll table. Return with Esc)")
	hintText.SetTextColor(tcell.ColorGray)
	body.AddItem(hintText, 4, 0, false)

	body.AddItem(primitives.TitleTextView("Balance"), 3, 0, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return primitives.NewCenterAlignedTextView(err.Error())
	}

	balanceTable := tview.NewTable().SetBorders(true)
	totalBalanceTextView := primitives.NewLeftAlignedTextView("")
	totalBalanceTextView.SetBorderPadding(0, 0, 5, 0)

	spendableBalanceTextView := primitives.NewLeftAlignedTextView("")
	spendableBalanceTextView.SetBorderPadding(0, 0, 5, 0)

	removeTextview := func(balanceTable, spendableBalanceTextView, totalBalanceTextView tview.Primitive) {
		body.RemoveItem(balanceTable)
		body.RemoveItem(spendableBalanceTextView)
		body.RemoveItem(totalBalanceTextView)
	}

	var simpleOutput, detailedOutput func([]*walletcore.Account)
	simpleOutput = func(accounts []*walletcore.Account) {
		//removeTextview clear the screen before outputing new data
		removeTextview(balanceTable, spendableBalanceTextView, totalBalanceTextView)

		if len(accounts) == 1 {
			output := walletcore.SimpleBalance(accounts[0].Balance, false)

			body.AddItem(spendableBalanceTextView.SetText(fmt.Sprintf("Total:        %s", output)), 2, 0, false)
			body.AddItem(totalBalanceTextView.SetText(fmt.Sprintf("Spendable:    %s", accounts[0].Balance.Spendable)), 2, 0, false)
		} else {
			balanceTable.SetCell(0, 0, tview.NewTableCell("Account Name").SetAlign(tview.AlignCenter)).
				SetCell(0, 1, tview.NewTableCell("Balance").SetAlign(tview.AlignCenter))

			for i, account := range accounts {
				row := i + 1
				balanceTable.SetCell(row, 0, tview.NewTableCell(account.Name).SetAlign(tview.AlignCenter)).
					SetCell(row, 1, tview.NewTableCell(walletcore.SimpleBalance(account.Balance, true)).SetAlign(tview.AlignCenter)).
					SetCell(row, 2, tview.NewTableCell(account.Balance.Spendable.String()).SetAlign(tview.AlignCenter))
			}

			body.AddItem(balanceTable, 0, 2, true)
		}
	}

	detailedOutput = func(accounts []*walletcore.Account) {
		//removeTextview clear the screen before outputing new data
		removeTextview(balanceTable, spendableBalanceTextView, totalBalanceTextView)

		balanceTable.SetCell(0, 0, tview.NewTableCell("Account Name").SetAlign(tview.AlignCenter)).
			SetCell(0, 1, tview.NewTableCell("Balance").SetAlign(tview.AlignCenter)).
			SetCell(0, 2, tview.NewTableCell("Spendable").SetAlign(tview.AlignCenter)).
			SetCell(0, 3, tview.NewTableCell("Locked").SetAlign(tview.AlignCenter)).
			SetCell(0, 4, tview.NewTableCell("Voting Authority").SetAlign(tview.AlignCenter)).
			SetCell(0, 5, tview.NewTableCell("Unconfirmed").SetAlign(tview.AlignCenter))

		for i, account := range accounts {
			row := i + 1
			balanceTable.SetCell(row, 0, tview.NewTableCell(account.Name).SetAlign(tview.AlignCenter)).
				SetCell(row, 1, tview.NewTableCell(walletcore.SimpleBalance(account.Balance, true)).SetAlign(tview.AlignCenter)).
				SetCell(row, 2, tview.NewTableCell(account.Balance.Spendable.String()).SetAlign(tview.AlignCenter)).
				SetCell(row, 3, tview.NewTableCell(account.Balance.LockedByTickets.String()).SetAlign(tview.AlignCenter)).
				SetCell(row, 4, tview.NewTableCell(account.Balance.VotingAuthority.String()).SetAlign(tview.AlignCenter)).
				SetCell(row, 5, tview.NewTableCell(account.Balance.Unconfirmed.String()).SetAlign(tview.AlignCenter))
		}

		body.AddItem(balanceTable, 0, 2, true)
	}

	//form button was used because tview buttons canot be style when part of flexbox
	detailedButtonForm := tview.NewForm()
	detailedButtonForm.SetBorderPadding(0, 0, 0, 0)
	body.AddItem(detailedButtonForm.AddButton("Detailed Balance", func() {
		detailedOutput(accounts)
	}), 3, 1, true)

	simpleOutput(accounts)

	// use different key press listener on first form item to watch for escape and enter key to set focus on table
	detailedButtonForm.GetButton(0).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			detailedOutput(accounts)
			setFocus(balanceTable)
			return nil
		}

		return event
	})

	// use different key press listener on first form item to watch for escape and enter key to set focus on button
	balanceTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			simpleOutput(accounts)
			setFocus(detailedButtonForm.GetButton(0))

			return nil
		}

		return event
	})

	setFocus(body)
	return body
}

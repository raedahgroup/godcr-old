package pages

import (
	"github.com/gdamore/tcell"
		"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func balancePage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	textView := tview.NewTextView()

	body := tview.NewFlex().SetDirection(tview.FlexRow)

	hintText := primitives.WordWrappedTextView("(TIP: Hit ENTER to view detailed balance, Navigate with Tab and Shift+Tab. Return with Esc)")
	hintText.SetTextColor(tcell.ColorGray)
	body.AddItem(hintText, 3, 0, false)

	body.AddItem(primitives.NewLeftAlignedTextView("Wallet Balance").SetTextColor(helpers.TitleColor), 3, 0, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return primitives.NewCenterAlignedTextView(err.Error())
	}

	checkbox := tview.NewCheckbox()
	table := tview.NewTable().SetBorders(true)
	textView := tview.NewTextView()
	textView.SetBorderPadding(0, 0, 30, 0)

	output := func (accounts []*walletcore.Account) {
		body.RemoveItem(table)

		if len(accounts) == 1 {
			output := walletcore.SimpleBalance(accounts[0].Balance, false)
			body.AddItem(textView.SetTextAlign(tview.AlignLeft).SetText(output), 5, 1, false)
		} else {
			for _, account := range accounts {
				body.AddItem(
					table.SetCell(0, 0, tview.NewTableCell("Account Name").SetAlign(tview.AlignCenter)).
						SetCell(0, 1, tview.NewTableCell("Balance").SetAlign(tview.AlignCenter)).
						SetCell(1, 0, tview.NewTableCell(account.Name).SetAlign(tview.AlignCenter)).
						SetCell(1, 1, tview.NewTableCell(walletcore.SimpleBalance(account.Balance, true)).SetAlign(tview.AlignCenter)), 0, 2, false)
			}
		}
	}

	body.AddItem(checkbox.SetLabel("View Detailed Balance  ").SetChecked(false).SetChangedFunc(func(checked bool) {
		if checked {
			body.RemoveItem(table)
			body.RemoveItem(textView)

			for _, account := range accounts {
				body.AddItem(
					table.SetCell(0, 0, tview.NewTableCell("Account Name").SetAlign(tview.AlignCenter)).
						SetCell(0, 1, tview.NewTableCell("Balance").SetAlign(tview.AlignCenter)).
						SetCell(0, 2, tview.NewTableCell("Spendable").SetAlign(tview.AlignCenter)).
						SetCell(0, 3, tview.NewTableCell("Locked").SetAlign(tview.AlignCenter)).
						SetCell(0, 4, tview.NewTableCell("Voting Authority").SetAlign(tview.AlignCenter)).
						SetCell(0, 5, tview.NewTableCell("Unconfirmed").SetAlign(tview.AlignCenter)).
						SetCell(1, 0, tview.NewTableCell(account.Name).SetAlign(tview.AlignCenter)).
						SetCell(1, 1, tview.NewTableCell(walletcore.SimpleBalance(account.Balance, true)).SetAlign(tview.AlignCenter)).
						SetCell(1, 2, tview.NewTableCell(account.Balance.Spendable.String()).SetAlign(tview.AlignCenter)).
						SetCell(1, 3, tview.NewTableCell(account.Balance.LockedByTickets.String()).SetAlign(tview.AlignCenter)).
						SetCell(1, 4, tview.NewTableCell(account.Balance.VotingAuthority.String()).SetAlign(tview.AlignCenter)).
						SetCell(1, 5, tview.NewTableCell(account.Balance.Unconfirmed.String()).SetAlign(tview.AlignCenter)), 0, 2, true)
			}
		} else {
			output(accounts)
		}
	}), 3, 1, true)	

	output(accounts)



	checkbox.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			checkbox.SetChecked(false)
			clearFocus()
		}
	}).SetLabelColor(tcell.ColorWhite)

	checkbox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyTAB {
			setFocus(table)
			return nil
		}

		return event
	})

	// use different key press listener on first form item to watch for backtab press and restore focus to stake info
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyTAB{
			setFocus(checkbox)
			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			setFocus(checkbox)
			return nil
		}
		return event
	})

	body.SetBorderPadding(1, 0, 1, 0)
	
	setFocus(body)
	return body
}

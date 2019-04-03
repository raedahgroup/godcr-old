package pages

import (
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func balancePage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	balancePage := tview.NewFlex().SetDirection(tview.FlexRow)
	balancePage.SetBorderPadding(1, 0, 1, 0)

	titleTextView := primitives.TitleTextView("")
	balancePage.AddItem(titleTextView.SetText("Balance"), 1, 0, false)

	hintText := primitives.WordWrappedTextView("")
	hintText.SetTextColor(tcell.ColorGray)
	balancePage.AddItem(hintText.SetText("(TIP: Press TAB to display detailed balance info, ESC to return to Navigation Menu.)"), 3, 0, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return primitives.NewCenterAlignedTextView(err.Error())
	}

	multipleAccountBalanceTable := tview.NewTable().SetBorders(true)
	detailedBalanceTable := tview.NewTable().SetBorders(true)
	singleAccountBalanceTable := tview.NewTable()

	// clearBalanceViews clear the screen before outputing new data
	clearBalanceViews := func(hintText, singleAccountBalanceTable, detailedBalanceTable, titleTextView, multipleAccountBalanceTable tview.Primitive) {
		balancePage.RemoveItem(detailedBalanceTable)
		balancePage.RemoveItem(multipleAccountBalanceTable)
		balancePage.RemoveItem(titleTextView)
		balancePage.RemoveItem(hintText)
		balancePage.RemoveItem(singleAccountBalanceTable)

	}

	var simpleOutput, detailedOutput func()
	simpleOutput = func() {
		clearBalanceViews(hintText, singleAccountBalanceTable, detailedBalanceTable, titleTextView, multipleAccountBalanceTable)

		balancePage.AddItem(titleTextView.SetText("Balance"), 1, 0, false)
		balancePage.AddItem(hintText.SetText("(TIP: Press TAB to display Detailed balance info, ESC to return to Navigation Menu.)"), 3, 0, false)

		if len(accounts) == 1 {
			output := walletcore.SimpleBalance(accounts[0].Balance, false)

			if accounts[0].Balance.Total == accounts[0].Balance.Spendable {
				singleAccountBalanceTable.SetCell(0, 0, tview.NewTableCell("Total: ").SetAlign(tview.AlignLeft)).
					SetCell(0, 1, tview.NewTableCell(output).SetAlign(tview.AlignRight))

				balancePage.AddItem(singleAccountBalanceTable, 0, 1, false)
			} else {
				singleAccountBalanceTable.SetCell(0, 0, tview.NewTableCell("Total: ").SetAlign(tview.AlignLeft)).
					SetCell(0, 1, tview.NewTableCell(output).SetAlign(tview.AlignRight)).
					SetCell(1, 0, tview.NewTableCell("Spendable: ").SetAlign(tview.AlignLeft)).
					SetCell(1, 1, tview.NewTableCell(accounts[0].Balance.Spendable.String()).SetAlign(tview.AlignRight))

				balancePage.AddItem(singleAccountBalanceTable, 0, 1, false)
			}
		} else {
			multipleAccountBalanceTable.SetCell(0, 0, tview.NewTableCell("Account Name").SetAlign(tview.AlignCenter)).
				SetCell(0, 1, tview.NewTableCell("Balance").SetAlign(tview.AlignCenter))

			for i, account := range accounts {
				row := i + 1
				if account.Balance.Total != account.Balance.Spendable {
					multipleAccountBalanceTable.SetCell(row, 0, tview.NewTableCell(account.Name).SetAlign(tview.AlignCenter)).
						SetCell(row, 1, tview.NewTableCell(walletcore.SimpleBalance(account.Balance, true)).SetAlign(tview.AlignCenter))

					balancePage.AddItem(multipleAccountBalanceTable, 0, 1, true)

				} else {
					multipleAccountBalanceTable.SetCell(0, 2, tview.NewTableCell("Spendable").SetAlign(tview.AlignCenter))
					multipleAccountBalanceTable.SetCell(row, 0, tview.NewTableCell(account.Name).SetAlign(tview.AlignCenter)).
						SetCell(row, 1, tview.NewTableCell(walletcore.SimpleBalance(account.Balance, true)).SetAlign(tview.AlignCenter)).
						SetCell(row, 2, tview.NewTableCell(account.Balance.Spendable.String()).SetAlign(tview.AlignCenter))

					balancePage.AddItem(multipleAccountBalanceTable, 0, 1, true)
				}
			}
		}

		// use tab and escape key press listener on simple balance when user has one account
		balancePage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				clearFocus()
				return nil
			}
			if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
				detailedOutput()
				setFocus(detailedBalanceTable)
				return nil
			}

			return event
		})

		// use tab and escape key press listener on simple balance table when user has more than one account
		multipleAccountBalanceTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				clearFocus()
				return nil
			}
			if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
				detailedOutput()
				setFocus(detailedBalanceTable)
				return nil
			}

			return event
		})
	}

	detailedOutput = func() {
		clearBalanceViews(hintText, singleAccountBalanceTable, detailedBalanceTable, titleTextView, multipleAccountBalanceTable)

		balancePage.AddItem(titleTextView.SetText("Balance (Detailed)"), 1, 0, false)
		balancePage.AddItem(hintText.SetText("(TIP: Press TAB to display Simple balance info, ARROW keys to Scroll table, ESC to return to Navigation Menu.)"), 3, 0, false)

		detailedBalanceTable.SetCell(0, 0, tview.NewTableCell("Account Name").SetAlign(tview.AlignCenter)).
			SetCell(0, 1, tview.NewTableCell("Balance").SetAlign(tview.AlignCenter)).
			SetCell(0, 2, tview.NewTableCell("Spendable").SetAlign(tview.AlignCenter)).
			SetCell(0, 3, tview.NewTableCell("Locked").SetAlign(tview.AlignCenter)).
			SetCell(0, 4, tview.NewTableCell("Voting Authority").SetAlign(tview.AlignCenter)).
			SetCell(0, 5, tview.NewTableCell("Unconfirmed").SetAlign(tview.AlignCenter))

		for i, account := range accounts {
			row := i + 1
			detailedBalanceTable.SetCell(row, 0, tview.NewTableCell(account.Name).SetAlign(tview.AlignCenter)).
				SetCell(row, 1, tview.NewTableCell(walletcore.SimpleBalance(account.Balance, true)).SetAlign(tview.AlignCenter)).
				SetCell(row, 2, tview.NewTableCell(account.Balance.Spendable.String()).SetAlign(tview.AlignCenter)).
				SetCell(row, 3, tview.NewTableCell(account.Balance.LockedByTickets.String()).SetAlign(tview.AlignCenter)).
				SetCell(row, 4, tview.NewTableCell(account.Balance.VotingAuthority.String()).SetAlign(tview.AlignCenter)).
				SetCell(row, 5, tview.NewTableCell(account.Balance.Unconfirmed.String()).SetAlign(tview.AlignCenter))
		}

		balancePage.AddItem(detailedBalanceTable, 0, 1, true)

		// use tab and escape key press listener on detailed balance table
		detailedBalanceTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				clearFocus()
				return nil
			}
			if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
				simpleOutput()
				setFocus(balancePage)
				return nil
			}

			return event
		})
	}

	simpleOutput()

	setFocus(balancePage)
	return balancePage
}

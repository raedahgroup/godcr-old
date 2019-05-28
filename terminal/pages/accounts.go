package pages

import (
"fmt"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func accountsPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	body.AddItem(primitives.NewLeftAlignedTextView("Accounts"), 2, 1, false)

	accountsTable := tview.NewTable().
		SetBorders(false).
		SetFixed(1, 0).
		SetSelectable(true, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		errorTextView := primitives.NewCenterAlignedTextView(err.Error()).SetTextColor(helpers.DecredOrangeColor)
		body.AddItem(errorTextView, 2, 0, false)
	}

	for _, account := range accounts {
		nextRowIndex := accountsTable.GetRowCount()

		accountsTable.SetCell(nextRowIndex, 0, tview.NewTableCell(fmt.Sprintf("%-10s", account.Name)).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1))
		accountsTable.SetCell(nextRowIndex, 1, tview.NewTableCell(fmt.Sprintf("%-10s", account.Balance.Total)).SetAlign(tview.AlignCenter).SetMaxWidth(2).SetExpansion(1))
		if account.Balance.Total != account.Balance.Spendable {
			accountsTable.SetCell(nextRowIndex, 2, tview.NewTableCell(fmt.Sprintf("%15s", account.Balance.Spendable)).SetAlign(tview.AlignCenter).SetMaxWidth(3).SetExpansion(1))
		}
	}

	body.AddItem(accountsTable, 0, 1, true)

	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			clearFocus()
			return nil
		}

		return event
	})

	setFocus(body)
	return body
}

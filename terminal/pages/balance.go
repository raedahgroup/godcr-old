package pages

import (
	"github.com/rivo/tview"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func BalancePage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	textView := tview.NewTextView()
	body := tview.NewFlex().SetDirection(tview.FlexRow)
	table := tview.NewTable().SetBorders(true)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return textView.SetText(err.Error())
	}
	
	body.AddItem(tview.NewCheckbox().SetLabel("View Detailed Balance  ").SetChecked(false).SetChangedFunc(func(checked bool) {
		if checked == false && len(accounts) == 1 {
			output := walletcore.SimpleBalance(accounts[0].Balance, false)
			body.AddItem(textView.SetTextAlign(tview.AlignCenter).SetText(output), 1, 1, false)
		}
		if checked == false && len(accounts) != 1 {
			for _, account := range accounts {
			body.AddItem(
				table.SetCell(0, 0, tview.NewTableCell("Account Name").SetAlign(tview.AlignCenter)).
				SetCell(0, 1, tview.NewTableCell("Balance").SetAlign(tview.AlignCenter)).
				SetCell(1, 0, tview.NewTableCell(account.Name).SetAlign(tview.AlignCenter)).
				SetCell(1, 1, tview.NewTableCell(walletcore.SimpleBalance(account.Balance, false)).SetAlign(tview.AlignCenter)), 1, 1, false)}
		}

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
			SetCell(1, 5, tview.NewTableCell(account.Balance.Unconfirmed.String()).SetAlign(tview.AlignCenter)), 0, 2, false)}
		}), 0, 2, true)
	
	setFocus(body)
	return body
}

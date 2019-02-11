package pages

import (
	"fmt"

	"github.com/rivo/tview"
	// "github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func BalancePage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewTextView().SetTextAlign(tview.AlignCenter)

	checkbox := tview.NewCheckbox().SetLabel("Show detailed ").SetChecked(false).Draw(screen)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		errMsg := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf(err.Error()))
		return errMsg
	}

	if len(accounts) == 1 {
		account := walletcore.SimpleBalance(accounts[0].Balance, false)
		body.SetText(account)
	}

	return body

	return nil
}

// func showSimpleView(window *nucular.Window) {
// 	helpers.SetFont(window, helpers.PageHeaderFont)
// 	window.Row(25).Dynamic(1)
// 	window.Label(walletcore.SimpleBalance(handler.accounts[0].Balance, false), "LC")
// }
package pages

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func BalancePage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	type body struct	{
		body []tview.Primitive
	} 

	body = tview.NewTextView().SetTextAlign(tview.AlignCenter)

	// body := tview.NewCheckbox().SetLabel("Show detailed ").SetChecked(false).Draw(screen)

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


	// accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	// if err != nil {
	// 	errMsg := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf(err.Error()))
	// 	return errMsg
	// }

	// body := tview.NewCheckbox().SetLabel("Show detailed ").SetChecked(false).SetChangedFunc(func(checked bool) {
	// 	if checked != true && len(accounts) == 1 {
	// 		account := walletcore.SimpleBalance(accounts[0].Balance, false)
	// 		 fmt.Println(account)
	// 	}else{
	// 		 fmt.Println("more than 1")
	// 	}
	// })
	// body.SetDoneFunc(func (key tcell.Key) {
	// 	if key == tcell.KeyEscape{
	// 		clearFocus()
	// 	}
	// })
	// setFocus(body)
	// return body
}

package pages

import (
	"github.com/raedahgroup/godcr/app"
	"github.com/rivo/tview"
)

func exitPage(walletMiddleware app.WalletMiddleware, tviewApp *tview.Application, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewModal().
		SetText("Do you want to quit Terminal application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				tviewApp.Stop()
			} else {
				tviewApp.SetRoot(rootPage(tviewApp, walletMiddleware), true)
			}
		})

	tviewApp.SetRoot(body, true)
	return body
}

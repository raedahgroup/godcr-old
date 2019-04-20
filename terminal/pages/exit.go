package pages

import "github.com/rivo/tview"

func exitPage(tviewApp *tview.Application, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewModal().
		SetText("Do you want to quit Terminal application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				tviewApp.Stop()
			} else {
				clearFocus()
			}
		})

	setFocus(body)
	return body
}

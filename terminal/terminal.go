package terminal

import (
	"github.com/rivo/tview"
)

func OpenTerminal() {

	app := tview.NewApplication()

	list := tview.NewList().
		AddItem("Open wallet", "", 'a', nil).
		AddItem("Sync on load", "", 'b', nil).
		AddItem("Account Balance", "", 'c', nil).
		AddItem("Close / Shutdown", "Press to exit", 'q', func() {
			app.Stop()
		})

	list.SetBorder(true).SetTitle("GoDCR Terminal")

	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}

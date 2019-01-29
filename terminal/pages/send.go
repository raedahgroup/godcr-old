package pages

import (
	"fmt"
	
	"github.com/rivo/tview"
)

func SendPage(tviewApp *tview.Application, menuColumn *tview.List) tview.Primitive {
	title := pageTitle("Send")

	//Form for Sending
	form := tview.NewForm().
	AddDropDown("Account", []string{"Dafault", "..."}, 0, nil).
	AddInputField("Amount", "", 20, nil, nil).
	AddInputField("Destination Address", "", 20, nil, nil).
	AddButton("Send", func() {
		fmt.Println("Next")
	})
	
	form.AddButton("Close", func() {
		tviewApp.SetFocus(menuColumn)
	})

	gridSend := tview.NewGrid().SetRows(2, 0).SetColumns(0)
	gridSend.AddItem(title, 0, 0, 1, 1, 0, 0, true).
				AddItem(form, 1, 0, 1, 1, 0, 0, true)

	tviewApp.SetFocus(form)
	return gridSend
}
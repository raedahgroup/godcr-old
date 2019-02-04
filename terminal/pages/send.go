package pages

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func SendPage(setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	//Form for Sending
	body := tview.NewForm().
		AddDropDown("Account", []string{"Dafault", "..."}, 0, nil).
		AddInputField("Amount", "", 20, nil, nil).
		AddInputField("Destination Address", "", 20, nil, nil).
		AddButton("Send", func() {
			fmt.Println("Next")
		})
	body.AddButton("Cancel", func() {
		clearFocus()
	})
	body.SetLabelColor(tcell.NewRGBColor(255, 255, 255))
	setFocus(body)

	return body
}

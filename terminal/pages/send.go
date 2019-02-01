package pages

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
)

var LabelColor, BorderColor tcell.Color = tcell.NewRGBColor(255, 255, 255), tcell.NewRGBColor(112, 203, 255)

func SendPage(setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	//Form for Sending
	body := tview.NewForm().
		AddDropDown("Account", []string{"Default", "..."}, 0, nil).
		AddInputField("Amount", "", 20, nil, nil).
		AddInputField("Destination Address", "", 20, nil, nil).
		AddButton("Send", func() {
			fmt.Println("Next")
		})
	body.AddButton("Cancel", func() {
		clearFocus()
	})
	body.SetBackgroundColor(tcell.NewRGBColor(255, 255, 255))
	body.SetLabelColor(tcell.NewRGBColor(0, 0, 0))
	body.SetFieldTextColor(tcell.NewRGBColor(0, 0, 0))

	body.SetLabelColor(LabelColor)
	setFocus(body)

	return body
}

package pages

import (
	"fmt"
	
	"github.com/rivo/tview"
)

func ReceivePage() tview.Primitive {
	title := pageTitle("Receive")

	body := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("ID : %s","TsU1YvSmtqw7wUtsRvSWjVs9BRxfT7urLzN"))
	
	gridReceive := tview.NewGrid().SetRows(2, 0).SetColumns(0)
	gridReceive.AddItem(title, 0, 0, 1, 1, 0, 0, true).SetBorders(false)
	gridReceive.AddItem(body, 1, 0, 1, 1, 0, 0, true).SetBorders(false)

	return gridReceive
}
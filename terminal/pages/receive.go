package pages

import (
	"fmt"
	
	"github.com/rivo/tview"
)

func ReceivePage() tview.Primitive {

	body := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("ID : %s","TsU1YvSmtqw7wUtsRvSWjVs9BRxfT7urLzN"))

	return body
}
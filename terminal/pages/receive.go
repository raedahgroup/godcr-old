package pages

import (
	"fmt"

	"github.com/rivo/tview"
)

func ReceivePage() tview.Primitive {
	body := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("\n\n\nID : %s", "TsU1YvSmtqw7wUtsRvSWjVs9BRxfT7urLzN"))
	body.SetTextColor(LabelColor)

	return body
}

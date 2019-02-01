package pages

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
)

func ReceivePage() tview.Primitive {
	body := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("\n\n\nID : %s", "TsU1YvSmtqw7wUtsRvSWjVs9BRxfT7urLzN"))
	body.SetTextColor(tcell.NewRGBColor(0, 0, 0))
	body.SetBackgroundColor(tcell.NewRGBColor(255, 255, 255))

	return body
}

package pages

import (
	"fmt"

	"github.com/rivo/tview"
)

func BalancePage() tview.Primitive {
	body := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("\n\n\nBalance : %s", "0 GODCR"))
	body.SetTextColor(tcell.NewRGBColor(0, 0, 0))
	body.SetBackgroundColor(tcell.NewRGBColor(255, 255, 255))

	return body
}

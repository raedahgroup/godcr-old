package pages

import (
	"fmt"

	"github.com/rivo/tview"
)

func BalancePage() tview.Primitive {
	body := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Balance : %s", "0 GODCR"))

	return body
}
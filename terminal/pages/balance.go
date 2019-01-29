package pages

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)


func BalancePage() tview.Primitive {
	title := pageTitle("Balance")
	body := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Balance : %s","0 GODCR"))
	gridBalance := tview.NewGrid().SetRows(2, 0).SetColumns(0)
	gridBalance.AddItem(title, 0, 0, 1, 1, 0, 0, true).SetBorders(false)
	gridBalance.AddItem(body, 1, 0, 1, 1, 0, 0, true).SetBorders(false)

	return gridBalance
}

func pageTitle(text string) tview.Primitive {
	return tview.NewTextView().SetTextAlign(tview.AlignLeft).SetText(fmt.Sprintf("Page::%s", strings.ToUpper(text)))
}
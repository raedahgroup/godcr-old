package pages

import (
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func accountsPage(setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	body.AddItem(primitives.NewLeftAlignedTextView("Accounts"), 2, 1, false)

	body.AddItem(primitives.NewLeftAlignedTextView("Accounts page coming soon").SetTextColor(helpers.DecredGreenColor), 0, 1, false)

	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			clearFocus()
			return nil
		}

		return event
	})

	setFocus(body)
	return body
}

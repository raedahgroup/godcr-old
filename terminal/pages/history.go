package pages

import (
	"github.com/rivo/tview"
)

func HistoryPage() tview.Primitive {
	title := pageTitle("History")

	body := tview.NewTable().SetBorders(true)
	body.SetCell(0, 0, tview.NewTableCell("Date").SetAlign(tview.AlignCenter))
	body.SetCell(0, 1, tview.NewTableCell("Amount").SetAlign(tview.AlignCenter))
	body.SetCell(0, 2, tview.NewTableCell("Fee").SetAlign(tview.AlignCenter))
	body.SetCell(0, 3, tview.NewTableCell("Direction").SetAlign(tview.AlignCenter))
	body.SetCell(0, 4, tview.NewTableCell("Type").SetAlign(tview.AlignCenter))

	gridReceive := tview.NewGrid().SetRows(2, 0).SetColumns(0)
	gridReceive.AddItem(title, 0, 0, 1, 1, 0, 0, true).SetBorders(false)
	gridReceive.AddItem(body, 1, 0, 1, 1, 0, 0, true).SetBorders(false)

	return gridReceive
}
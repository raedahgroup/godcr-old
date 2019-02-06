package pages

import (
	"github.com/rivo/tview"
)

func HistoryPage() tview.Primitive {
	body := tview.NewTable().SetBorders(true)
	body.SetCell(0, 0, tview.NewTableCell("Date").SetAlign(tview.AlignCenter).SetTextColor(LabelColor))
	body.SetCell(0, 1, tview.NewTableCell("Amount").SetAlign(tview.AlignCenter).SetTextColor(LabelColor))
	body.SetCell(0, 2, tview.NewTableCell("Fee").SetAlign(tview.AlignCenter).SetTextColor(LabelColor))
	body.SetCell(0, 3, tview.NewTableCell("Direction").SetAlign(tview.AlignCenter).SetTextColor(LabelColor))
	body.SetCell(0, 4, tview.NewTableCell("Type").SetAlign(tview.AlignCenter).SetTextColor(LabelColor))

	return body
}

package pages

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func HistoryPage() tview.Primitive {
	color := tcell.NewRGBColor(255, 255, 255)
	body := tview.NewTable().SetBorders(true)
	body.SetCell(0, 0, tview.NewTableCell("Date").SetAlign(tview.AlignCenter).SetTextColor(color))
	body.SetCell(0, 1, tview.NewTableCell("Amount").SetAlign(tview.AlignCenter).SetTextColor(color))
	body.SetCell(0, 2, tview.NewTableCell("Fee").SetAlign(tview.AlignCenter).SetTextColor(color))
	body.SetCell(0, 3, tview.NewTableCell("Direction").SetAlign(tview.AlignCenter).SetTextColor(color))
	body.SetCell(0, 4, tview.NewTableCell("Type").SetAlign(tview.AlignCenter).SetTextColor(color))
	body.SetBackgroundColor(tcell.NewRGBColor(255, 255, 255))
	body.SetBorderColor(tcell.NewRGBColor(0, 0, 51))
	body.SetBorderAttributes(tcell.AttrBold)

	return body
}

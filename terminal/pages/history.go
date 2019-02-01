package pages

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
)

func HistoryPage() tview.Primitive {
	body := tview.NewTable().SetBorders(true)
	body.SetCell(0, 0, tview.NewTableCell("Date").SetAlign(tview.AlignCenter))
	body.SetCell(0, 1, tview.NewTableCell("Amount").SetAlign(tview.AlignCenter))
	body.SetCell(0, 2, tview.NewTableCell("Fee").SetAlign(tview.AlignCenter))
	body.SetCell(0, 3, tview.NewTableCell("Direction").SetAlign(tview.AlignCenter))
	body.SetCell(0, 4, tview.NewTableCell("Type").SetAlign(tview.AlignCenter))
	body.SetBackgroundColor(tcell.NewRGBColor(255, 255, 255))
	body.SetBorderColor(tcell.NewRGBColor(0, 0, 51))
	body.SetBorderAttributes(tcell.AttrBold)

	return body
}
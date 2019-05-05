package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func HistoryPage(windows fyne.Window, App fyne.App) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("History", fyne.TextAlignLeading, fyne.TextStyle{Italic: true, Bold: true})
	table := widgets.NewTable()
	table, _ = FetchRecentActivity(Wallet, table, 0, true)
	info := widget.NewVBox(
		label,
		widgets.NewVSpacer(30),
		table.CondensedTable())
	return widget.NewScrollContainer(info)
}

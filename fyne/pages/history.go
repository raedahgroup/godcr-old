package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func HistoryPage(offSet, count int, windows fyne.Window, App fyne.App) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("History", fyne.TextAlignLeading, fyne.TextStyle{Italic: true, Bold: true})
	table := widgets.NewTable()
	table, _ = FetchRecentActivity(Wallet, table, offSet, count, true)

	totalNoOfTxns, _ := Wallet.TransactionCount(nil)

	var next *widget.Button
	next = widget.NewButton("Next", func() {
		windows.SetContent(Menu(widget.NewLabelWithStyle("fetching data...", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true, Bold: true}), windows, App))
		widget.Refresh(next)
		if totalNoOfTxns-(offSet+15) > 0 {
			windows.SetContent(Menu(HistoryPage(offSet+count, 15, windows, App), windows, App))
		} else if totalNoOfTxns-offSet > 0 {
			windows.SetContent(Menu(HistoryPage(offSet+count, totalNoOfTxns-offSet, windows, App), windows, App))
		}
	})
	var back *widget.Button
	back = widget.NewButton("Back", func() {
		windows.SetContent(Menu(widget.NewLabelWithStyle("fetching data...", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true, Bold: true}), windows, App))
		if offSet-15 <= 15 {
			windows.SetContent(Menu(HistoryPage(1, 15, windows, App), windows, App))
		} else {
			windows.SetContent(Menu(HistoryPage(offSet-15, 15, windows, App), windows, App))
		}
	})

	if offSet == 1 {
		widget.Refresh(back)
		back.Hide()
	}
	if offSet+15 >= totalNoOfTxns {
		widget.Refresh(next)
		next.Hide()
	}

	buttons := widget.NewHBox(
		back,
		widgets.NewHSpacer(10),
		next,
	)

	return widget.NewVBox(
		label,
		widgets.NewVSpacer(10),
		table.CondensedTable(),
		widgets.NewVSpacer(10),
		buttons)
}

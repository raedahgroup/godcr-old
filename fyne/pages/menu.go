package pages

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

var (
	PeerConn  *canvas.Text
	BlkHeight *canvas.Text
)

func PageNotImplemented(win fyne.Window, App fyne.App) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	return Menu(label, win, App)
}

//Menu outputs CanvasObjects "widgets" in vertical formats
func Menu(object fyne.CanvasObject, win fyne.Window, App fyne.App) fyne.CanvasObject {

	PeerConn = canvas.NewText("", color.Black)
	BlkHeight = canvas.NewText("", color.Black)

	BlkHeight.TextStyle = fyne.TextStyle{Bold: true}
	PeerConn.TextStyle = fyne.TextStyle{Bold: true}
	BlkHeight.Alignment = fyne.TextAlignCenter
	PeerConn.Alignment = fyne.TextAlignCenter

	button := Button(App, win)
	menuSectionSize := fyne.NewSize(200, button.MinSize().Height)
	objContainer := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, widgets.NewHSpacer(20), widgets.NewHSpacer(20)), object)
	menuSize := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(menuSectionSize), button)
	box := widget.NewVBox(menuSize, widgets.NewVSpacer(10), widgets.NewHSpacer(100), PeerConn, widgets.NewVSpacer(10), widgets.NewHSpacer(100), BlkHeight)

	return fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, box, nil), box, objContainer)
}

func Button(App fyne.App, window fyne.Window) fyne.CanvasObject {
	return widget.NewGroup("Menu",

		widget.NewButton("Overview", func() {
			if SyncDone == true {
				window.SetContent(Menu(widget.NewLabelWithStyle("fetching data...", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true, Bold: true}), window, App))
				go window.SetContent(Menu(OverviewPage(window, App), window, App))
			}
		}),

		widget.NewButton("History", func() {
			if SyncDone == true {
				window.SetContent(Menu(widget.NewLabelWithStyle("fetching data...", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true, Bold: true}), window, App))
				window.SetContent(Menu(HistoryPage(window, App), window, App))
			}
		}),

		widget.NewButton("Send", func() {
			if SyncDone == true {
				window.SetContent(PageNotImplemented(window, App))
			}
		}),

		widget.NewButton("Receive", func() {
			if SyncDone == true {
				window.SetContent(PageNotImplemented(window, App))
			}
		}),

		widget.NewButton("Staking", func() {
			if SyncDone == true {
				window.SetContent(PageNotImplemented(window, App))
			}
		}),

		widget.NewButton("Accounts", func() {
			if SyncDone == true {
				window.SetContent(PageNotImplemented(window, App))
			}
		}),

		widget.NewButton("Security", func() {
			if SyncDone == true {
				window.SetContent(PageNotImplemented(window, App))
			}
		}),

		widget.NewButton("Settings", func() {
			if SyncDone == true {
				window.SetContent(SettingsPage(window, App))
			}
		}),

		widget.NewButton("Exit", func() {
			window.Close()
			App.Quit()
			//return
		}),
	)
}

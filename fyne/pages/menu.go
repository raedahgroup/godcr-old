package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
)

var (
	PeerConn  *widget.Label
	BlkHeight *widget.Label
)

func init() {
	PeerConn = widget.NewLabel("")
	BlkHeight = widget.NewLabel("")
}

func pageNotImplemented() fyne.CanvasObject {
	label := widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	return label
}

func Menu(wallet godcrApp.WalletMiddleware) fyne.CanvasObject {
	tabs := widget.NewTabContainer(
		widget.NewTabItem("    Overview", overviewPage(wallet)),
		widget.NewTabItem("    History", pageNotImplemented()),
		widget.NewTabItem("    Send", pageNotImplemented()),
		widget.NewTabItem("    Receive", pageNotImplemented()),
		widget.NewTabItem("    Staking", pageNotImplemented()),
		widget.NewTabItem("    Accounts", pageNotImplemented()),
		widget.NewTabItem("    Security", pageNotImplemented()),
		widget.NewTabItem("    Settings", pageNotImplemented()))

	tabs.SetTabLocation(widget.TabLocationLeading)
	text := layout.NewVBoxLayout()
	orderedText := []fyne.CanvasObject{PeerConn, BlkHeight}
	text.Layout(orderedText, fyne.NewSize(0, 0))
	textObj := fyne.NewContainerWithLayout(text, PeerConn, BlkHeight)
	tAlign := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, textObj, nil), textObj)
	menu := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, tAlign, tabs, nil), tabs, tAlign)
	return menu
}

// //Menu outputs CanvasObjects "widgets" in vertical formats
// func Menu(object fyne.CanvasObject, win fyne.Window, App fyne.App) fyne.CanvasObject {

// 	PeerConn = canvas.NewText("", color.Black)
// 	BlkHeight = canvas.NewText("", color.Black)

// 	BlkHeight.TextStyle = fyne.TextStyle{Bold: true}
// 	PeerConn.TextStyle = fyne.TextStyle{Bold: true}
// 	BlkHeight.Alignment = fyne.TextAlignCenter
// 	PeerConn.Alignment = fyne.TextAlignCenter

// 	button := Button(App, win)
// 	menuSectionSize := fyne.NewSize(200, button.MinSize().Height)
// 	objContainer := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, widgets.NewHSpacer(20), widgets.NewHSpacer(20)), object)
// 	menuSize := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(menuSectionSize), button)
// 	box := widget.NewVBox(menuSize, widgets.NewVSpacer(10), widgets.NewHSpacer(100), PeerConn, widgets.NewVSpacer(10), widgets.NewHSpacer(100), BlkHeight)

// 	return fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, box, nil), box, objContainer)
// }

// func Button(App fyne.App, window fyne.Window) fyne.CanvasObject {
// 	return widget.NewGroup("Menu",

// 		widget.NewButton("Overview", func() {
// 			if SyncDone == true {
// 				window.SetContent(Menu(widget.NewLabelWithStyle("fetching data...", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true, Bold: true}), window, App))
// 				window.SetContent(Menu(OverviewPage(window, App), window, App))
// 			}
// 		}),

// 		widget.NewButton("History", func() {
// 			if SyncDone == true {
// 				window.SetContent(Menu(widget.NewLabelWithStyle("fetching data...", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true, Bold: true}), window, App))
// 				window.SetContent(Menu(HistoryPage(1, 15, window, App), window, App))
// 			}
// 		}),

// 		widget.NewButton("Send", func() {
// 			if SyncDone == true {
// 				window.SetContent(PageNotImplemented(window, App))
// 			}
// 		}),

// 		widget.NewButton("Receive", func() {
// 			if SyncDone == true {
// 				window.SetContent(PageNotImplemented(window, App))
// 			}
// 		}),

// 		widget.NewButton("Staking", func() {
// 			if SyncDone == true {
// 				window.SetContent(PageNotImplemented(window, App))
// 			}
// 		}),

// 		widget.NewButton("Accounts", func() {
// 			if SyncDone == true {
// 				window.SetContent(PageNotImplemented(window, App))
// 			}
// 		}),

// 		widget.NewButton("Security", func() {
// 			if SyncDone == true {
// 				window.SetContent(PageNotImplemented(window, App))
// 			}
// 		}),

// 		widget.NewButton("Settings", func() {
// 			if SyncDone == true {
// 				window.SetContent(SettingsPage(window, App))
// 			}
// 		}),

// 		widget.NewButton("Exit", func() {
// 			window.Close()
// 			App.Quit()
// 			//return
// 		}),
// 	)
// }

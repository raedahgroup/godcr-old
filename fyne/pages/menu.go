package pages

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/canvas"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func PageNotImplemented(win fyne.Window, App fyne.App) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	return Menu(label, App, win)
}
func FecthingDataLabel(win fyne.Window, App fyne.App) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("fetching data...", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Bold: true})
	return Menu(label, App, win)
}

//Menu outputs CanvasObjects "widgets" in vertical formats
func Menu(object fyne.CanvasObject, App fyne.App, win fyne.Window) fyne.CanvasObject {
	var text2 *canvas.Text
	var text1 *canvas.Text

	text2 = canvas.NewText("", color.Black)
	text1 = canvas.NewText("", color.Black)

	text1.TextStyle = fyne.TextStyle{Bold: true}
	text2.TextStyle = fyne.TextStyle{Bold: true}
	text1.Alignment = fyne.TextAlignCenter
	text2.Alignment = fyne.TextAlignCenter

	go func() {
		for i := 0; i < 10; i++ {
			sec, _ := time.ParseDuration("2s")
			time.Sleep(sec)
			fmt.Println("n")
			text2.Text = "Helloo " + strconv.Itoa(i)
			text1.Text = "Hello " + strconv.Itoa(i)
			canvas.Refresh(text2)
			canvas.Refresh(text1)
		}
	}()

	button := Button(App, win)
	menuSectionSize := fyne.NewSize(200, button.MinSize().Height)
	objContainer := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, widgets.NewHSpacer(20), widgets.NewHSpacer(20)), object)
	menuSize := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(menuSectionSize), button)
	box := widget.NewVBox(menuSize, widgets.NewVSpacer(10), widgets.NewHSpacer(100), text1, widgets.NewVSpacer(10), widgets.NewHSpacer(100), text2)
	
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, box, nil), box, objContainer)
}

func Button(App fyne.App, window fyne.Window) fyne.CanvasObject {
	return widget.NewGroup("Menu",

		widget.NewButton("Overview", func() {
			if SyncDone == true {
				window.SetContent(PageNotImplemented(window, App))
			}
		}),

		widget.NewButton("History", func() {
			if SyncDone == true {
				window.SetContent(PageNotImplemented(window, App))
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
			App.Quit()
		}),
	)
}

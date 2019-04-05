package fyne

import (
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	godcrApp "github.com/raedahgroup/godcr/app"
)

func (app *fyneApp) showSyncWindow() {
	window := app.NewWindow(godcrApp.DisplayName)

	// todo create the sync window content (widgets to show sync progress)
	syncWindowContent := widget.NewVBox(
		widget.NewLabelWithStyle("Sync progress will appear here after sync is implemented.", fyne.TextAlignCenter, fyne.TextStyle{Italic:true}),
		widget.NewLabelWithStyle("App will launch fully in a few seconds.", fyne.TextAlignCenter, fyne.TextStyle{Italic:true}),
	)

	window.SetContent(syncWindowContent)
	window.CenterOnScreen()
	window.Show()

	// todo begin sync process, when sync completes, close this window and show the main window
	go func() {
		time.Sleep(5 * time.Second)
		window.Close()

		app.mainWindow.Show()
		app.menuButtons[0].OnTapped()
	}()

}

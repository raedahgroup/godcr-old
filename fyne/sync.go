package fyne

import (
	godcrApp "github.com/raedahgroup/godcr/app"
	"fyne.io/fyne/widget"
	"time"
)

func (app *fyneApp) showSyncWindow() {
	window := app.NewWindow(godcrApp.DisplayName)

	// todo create the sync window content (widgets to show sync progress) and begin sync process
	// when sync completes, show the main window and close this window
	go func() {
		time.Sleep(2 * time.Second)
		app.mainWindow.Show()
		window.Close()
	}()

	// todo this window should contain the sync progress indicator
	window.SetContent(widget.NewLabel("Sync is not implemented yet but this window will automatically close after a brief moment"))

	window.CenterOnScreen()
	window.Show()
}

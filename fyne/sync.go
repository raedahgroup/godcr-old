package fyne

import (
	"fmt"
	"math"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

const (
	defaultWindowWidth  = 800
	defaultWindowHeight = 600
)

func (app *fyneApp) showSyncWindow() {
	statusLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})

	progressBar := widget.NewProgressBar()
	progressBar.Max = 100
	progressBar.Min = 0
	progressBar.Hide()

	syncWindowContent := widget.NewVBox(
		widgets.NewVSpacer(10),
		widget.NewLabelWithStyle("Synchronizing", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
		statusLabel,
		progressBar,
		widgets.NewVSpacer(20),
	)

	syncWindowHeight := syncWindowContent.MinSize().Height + theme.Padding()*2
	syncWindowWidth := int(math.Round(defaultWindowWidth / 1.5))

	app.mainWindow.SetTitle(godcrApp.DisplayName)
	app.resizeAndCenterMainWindow(syncWindowContent, fyne.NewSize(syncWindowWidth, syncWindowHeight))
	app.mainWindow.Show()

	err := app.walletMiddleware.SyncBlockChainOld(&godcrApp.BlockChainSyncListener{
		SyncStarted: func() {
			statusLabel.SetText("Sync started...")
		},
		SyncEnded: func(err error) {
			progressBar.Hide()
			if err != nil {
				statusLabel.SetText(fmt.Sprintf("Sync completed with error: %s", err.Error()))
			} else {
				statusLabel.SetText("Sync completed successfully")
				time.Sleep(1 * time.Second)
				app.loadMainWindowContent()
			}
		},
		OnHeadersFetched: func(percentageProgress int64) {
			statusLabel.SetText(fmt.Sprintf("Blockchain sync in progress. Fetching headers (1/3): %d%%", percentageProgress))
			progressBar.Value = float64(percentageProgress)
			progressBar.Show()
		},
		OnDiscoveredAddress: func(_ string) {
			statusLabel.SetText("Blockchain sync in progress. Discovering addresses (2/3)")
			progressBar.Hide()
		},
		OnRescanningBlocks: func(percentageProgress int64) {
			statusLabel.SetText(fmt.Sprintf("Blockchain sync in progress. Rescanning blocks (3/3): %d%%", percentageProgress))
			progressBar.Value = float64(percentageProgress)
			progressBar.Show()
		},
	}, false)

	if err != nil {
		statusLabel.SetText(fmt.Sprintf("Sync failed to start: %s", err.Error()))
	}
}

func (app *fyneApp) loadMainWindowContent() {
	app.menuButtons[0].OnTapped()
	app.resizeAndCenterMainWindow(app.mainWindowContent, fyne.NewSize(defaultWindowWidth, defaultWindowHeight))
}

func (app *fyneApp) resizeAndCenterMainWindow(windowContent fyne.CanvasObject, size fyne.Size) {
	// create a fixedgrid wrapper around window content so that window.CenterOnScreen will work
	app.mainWindow.SetContent(fyne.NewContainerWithLayout(layout.NewFixedGridLayout(size), windowContent))
	app.mainWindow.CenterOnScreen()
}

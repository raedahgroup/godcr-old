package fyne

import (
	"fmt"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/sync"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

const (
	defaultWindowWidth  = 800
	defaultWindowHeight = 600
)

func (app *fyneApp) showSyncWindow() {
	progressBar := widget.NewProgressBar()
	progressBar.Max = 100
	progressBar.Min = 0

	reportLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	errorLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	var showDetails bool
	var fullSyncReport string
	var showDetailsButton *widget.Button
	showDetailsButton = widget.NewButton("Tap to view information", func() {
		showDetails = true
		reportLabel.SetText(fullSyncReport)
		showDetailsButton.Hide()
	})

	syncWindowContent := widget.NewVBox(
		widgets.NewVSpacer(10),
		widget.NewLabelWithStyle("Synchronizing", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
		progressBar,
		reportLabel,
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(showDetailsButton.MinSize()), showDetailsButton),
		errorLabel,
	)

	app.mainWindow.SetTitle(godcrApp.DisplayName)
	app.resizeAndCenterMainWindow(syncWindowContent)
	app.mainWindow.Show()

	err := app.walletMiddleware.SyncBlockChain(false, func(syncPrivateInfo *sync.PrivateInfo) {
		syncInfo := syncPrivateInfo.Read()
		progressBar.SetValue(float64(syncInfo.TotalSyncProgress))

		if syncInfo.Status == sync.StatusSuccess {
			app.loadMainWindowContent()
			return
		}

		report := strings.Builder{}
		if syncInfo.TotalTimeRemaining == "" {
			report.WriteString(fmt.Sprintf("%d%% completed.\n", syncInfo.TotalSyncProgress))
		} else {
			report.WriteString(fmt.Sprintf("%d%% completed, %s remaining.\n", syncInfo.TotalSyncProgress, syncInfo.TotalTimeRemaining))
		}

		if !showDetails {
			reportLabel.SetText(strings.TrimSpace(report.String()))
		}

		switch syncInfo.CurrentStep {
		case 1:
			report.WriteString(fmt.Sprintf("Fetched %d of %d block headers.\n", syncInfo.FetchedHeadersCount, syncInfo.TotalHeadersToFetch))
			report.WriteString(fmt.Sprintf("%d%% through step 1 of 3.\n", syncInfo.HeadersFetchProgress))

			if syncInfo.DaysBehind != "" {
				report.WriteString(fmt.Sprintf("Your wallet is %s behind.\n", syncInfo.DaysBehind))
			}

		case 2:
			report.WriteString("Discovering used addresses.\n")
			if syncInfo.AddressDiscoveryProgress > 100 {
				report.WriteString(fmt.Sprintf("%d%% (over) through step 2 of 3.\n", syncInfo.AddressDiscoveryProgress))
			} else {
				report.WriteString(fmt.Sprintf("%d%% through step 2 of 3.\n", syncInfo.AddressDiscoveryProgress))
			}

		case 3:
			report.WriteString(fmt.Sprintf("Scanning %d of %d block headers.\n", syncInfo.CurrentRescanHeight,
				syncInfo.TotalHeadersToFetch))
			report.WriteString(fmt.Sprintf("%d%% through step 3 of 3.\n", syncInfo.HeadersFetchProgress))
		}

		// show peer count last
		netType := app.walletMiddleware.NetType()
		if syncInfo.ConnectedPeers == 1 {
			report.WriteString(fmt.Sprintf("Syncing with %d peer on %s.\n", syncInfo.ConnectedPeers, netType))
		} else {
			report.WriteString(fmt.Sprintf("Syncing with %d peers on %s.\n", syncInfo.ConnectedPeers, netType))
		}

		fullSyncReport = report.String()
		if showDetails {
			reportLabel.SetText(fullSyncReport)
		}
	})
	
	if err != nil {
		errorMessage := fmt.Sprintf("Sync failed to start: %s", err.Error())
		errorLabel.SetText(errorMessage)
	} else {
		fullSyncReport = "Starting..."
		reportLabel.SetText("Starting...")
	}
}

func (app *fyneApp) loadMainWindowContent() {
	app.menuButtons[0].OnTapped()
	app.resizeAndCenterMainWindow(app.mainWindowContent)
}

func (app *fyneApp) resizeAndCenterMainWindow(windowContent fyne.CanvasObject) {
	// create a fixedgrid wrapper around window content so that window.CenterOnScreen will work
	windowSize := fyne.NewSize(defaultWindowWidth, defaultWindowHeight)
	app.mainWindow.SetContent(fyne.NewContainerWithLayout(layout.NewFixedGridLayout(windowSize), windowContent))
	app.mainWindow.CenterOnScreen()
}

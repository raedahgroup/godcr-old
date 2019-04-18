package fyne

import (
	"fmt"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
	godcrApp "github.com/raedahgroup/godcr/app"
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

	app.walletMiddleware.SyncBlockChain(false, func(report *defaultsynclistener.ProgressReport) {
		progressReport := report.Read()

		progressBar.SetValue(float64(progressReport.TotalSyncProgress))

		if progressReport.Status == defaultsynclistener.SyncStatusSuccess {
			app.loadMainWindowContent()
			return
		}

		stringReport := strings.Builder{}
		if progressReport.TotalTimeRemaining == "" {
			stringReport.WriteString(fmt.Sprintf("%d%% completed.\n", progressReport.TotalSyncProgress))
		} else {
			stringReport.WriteString(fmt.Sprintf("%d%% completed, %s remaining.\n", progressReport.TotalSyncProgress, progressReport.TotalTimeRemaining))
		}

		if !showDetails {
			reportLabel.SetText(strings.TrimSpace(stringReport.String()))
		}

		switch progressReport.CurrentStep {
		case defaultsynclistener.FetchingBlockHeaders:
			stringReport.WriteString(fmt.Sprintf("Fetched %d of %d block headers.\n", progressReport.FetchedHeadersCount, progressReport.TotalHeadersToFetch))
			stringReport.WriteString(fmt.Sprintf("%d%% through step 1 of 3.\n", progressReport.HeadersFetchProgress))

			if progressReport.DaysBehind != "" {
				stringReport.WriteString(fmt.Sprintf("Your wallet is %s behind.\n", progressReport.DaysBehind))
			}

		case defaultsynclistener.DiscoveringUsedAddresses:
			stringReport.WriteString("Discovering used addresses.\n")
			if progressReport.AddressDiscoveryProgress > 100 {
				stringReport.WriteString(fmt.Sprintf("%d%% (over) through step 2 of 3.\n", progressReport.AddressDiscoveryProgress))
			} else {
				stringReport.WriteString(fmt.Sprintf("%d%% through step 2 of 3.\n", progressReport.AddressDiscoveryProgress))
			}

		case defaultsynclistener.ScanningBlockHeaders:
			stringReport.WriteString(fmt.Sprintf("Scanning %d of %d block headers.\n", progressReport.CurrentRescanHeight,
				progressReport.TotalHeadersToFetch))
			stringReport.WriteString(fmt.Sprintf("%d%% through step 3 of 3.\n", progressReport.HeadersFetchProgress))
		}

		// show peer count last
		netType := app.walletMiddleware.NetType()
		if progressReport.ConnectedPeers == 1 {
			stringReport.WriteString(fmt.Sprintf("Syncing with %d peer on %s.\n", progressReport.ConnectedPeers, netType))
		} else {
			stringReport.WriteString(fmt.Sprintf("Syncing with %d peers on %s.\n", progressReport.ConnectedPeers, netType))
		}

		fullSyncReport = stringReport.String()
		if showDetails {
			reportLabel.SetText(fullSyncReport)
		}
	})
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

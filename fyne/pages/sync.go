package pages

import (
	"context"
	"fmt"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func ShowSyncWindow(ctx context.Context, wallet godcrApp.WalletMiddleware, window fyne.Window, App fyne.App) fyne.CanvasObject {
	progressBar := widget.NewProgressBar()
	progressBar.Min = 0
	progressBar.Max = 100

	var fullSyncReport string

	reportLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	widget.Refresh(reportLabel)
	reportLabel.Hide()
	var infoButton *widget.Button

	infoButton = widget.NewButton("Tap to view informations", func() {
		reportLabel.Show()
		infoButton.Hide()
	})

	var syncDone bool

	wallet.SyncBlockChain(false, func(report *defaultsynclistener.ProgressReport) {
		progressReport := report.Read()
		progressBar.SetValue(float64(progressReport.TotalSyncProgress))

		if progressReport.Status == defaultsynclistener.SyncStatusSuccess {
			if syncDone == false {
				syncDone = true
				menu := menuPage(ctx, wallet, App, window)
				window.SetContent(menu)
			}
		}

		stringReport := strings.Builder{}
		if progressReport.TotalTimeRemaining == "" {
			stringReport.WriteString(fmt.Sprintf("%d%% completed.\n", progressReport.TotalSyncProgress))
		} else {
			stringReport.WriteString(fmt.Sprintf("%d%% completed, %s remaining.\n", progressReport.TotalSyncProgress, progressReport.TotalTimeRemaining))
		}

		widget.Refresh(reportLabel)
		reportLabel.SetText(strings.TrimSpace(stringReport.String()))

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
		netType := wallet.NetType()
		if progressReport.ConnectedPeers == 1 {
			stringReport.WriteString(fmt.Sprintf("Syncing with %d peer on %s.\n", progressReport.ConnectedPeers, netType))
		} else {
			stringReport.WriteString(fmt.Sprintf("Syncing with %d peers on %s.\n", progressReport.ConnectedPeers, netType))
		}

		fullSyncReport = stringReport.String()

		widget.Refresh(reportLabel)
		reportLabel.SetText(fullSyncReport)
	})

	return widget.NewVBox(
		widgets.NewVSpacer(10),
		widget.NewLabelWithStyle("Synchronizing....", fyne.TextAlignLeading, fyne.TextStyle{Italic: true, Bold: true}),
		widgets.NewVSpacer(10),
		progressBar,
		widget.NewHBox(layout.NewSpacer(), infoButton, layout.NewSpacer()),
		reportLabel)
}

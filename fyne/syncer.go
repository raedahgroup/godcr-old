package fyne

import (
	"fmt"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type Syncer struct {
	err                error
	percentageProgress int
	report             []string
	showDetails        bool
	status             defaultsynclistener.SyncStatus
	syncError          error

	progressBar    *widget.ProgressBar
	reportLabel    *widget.Label
	errorLabel     *widget.Label
	fullSyncReport string
}

func NewSyncer() *Syncer {
	return &Syncer{
		percentageProgress: 0,
		syncError:          nil,
		report: []string{
			"Starting...",
		},
		showDetails: false,
	}
}

func (s *Syncer) startSyncing(walletMiddleware app.WalletMiddleware, changePageFunc func()) {
	//prog
	s.progressBar = widget.NewProgressBar()
	s.progressBar.Max = 100
	s.progressBar.Min = 0

	s.reportLabel = widget.NewLabel("")
	s.errorLabel = widget.NewLabel("")

	// begin block chain sync now so that when `Render` is called shortly after this, there'd be a report to display
	walletMiddleware.SyncBlockChain(false, func(report *defaultsynclistener.ProgressReport) {
		progressReport := report.Read()

		s.progressBar.SetValue(float64(progressReport.TotalSyncProgress))

		if progressReport.Status == defaultsynclistener.SyncStatusSuccess {
			s.status = defaultsynclistener.SyncStatusSuccess
			changePageFunc()
			return
		}

		stringReport := strings.Builder{}
		if progressReport.TotalTimeRemaining == "" {
			stringReport.WriteString(fmt.Sprintf("%d%% completed.\n", progressReport.TotalSyncProgress))
		} else {
			stringReport.WriteString(fmt.Sprintf("%d%% completed, %s remaining.\n", progressReport.TotalSyncProgress, progressReport.TotalTimeRemaining))
		}

		if !s.showDetails {
			s.reportLabel.SetText(strings.TrimSpace(stringReport.String()))
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
		netType := walletMiddleware.NetType()
		if progressReport.ConnectedPeers == 1 {
			stringReport.WriteString(fmt.Sprintf("Syncing with %d peer on %s.\n", progressReport.ConnectedPeers, netType))
		} else {
			stringReport.WriteString(fmt.Sprintf("Syncing with %d peers on %s.\n", progressReport.ConnectedPeers, netType))
		}

		s.fullSyncReport = stringReport.String()
		if s.showDetails {
			s.reportLabel.SetText(s.fullSyncReport)
		}
	})
}

func (s *Syncer) isDoneSyncing() bool {
	return s.status == defaultsynclistener.SyncStatusSuccess
}

func (s *Syncer) Render(container *widgets.Box) {
	var showDetailsButton *widget.Button
	showDetailsButton = widget.NewButton("Tap to view information", func() {
		s.showDetails = true
		s.reportLabel.SetText(s.fullSyncReport)
		showDetailsButton.Hide()
	})

	view := widgets.NewVBox(
		widgets.NewVSpacer(10),
	)
	view.AddItalicLabel("Synchronizing")
	view.Add(s.progressBar)
	view.Add(s.reportLabel)
	view.Add(fyne.NewContainerWithLayout(layout.NewFixedGridLayout(showDetailsButton.MinSize()), showDetailsButton))
	view.Add(s.errorLabel)

	container.Add(view)
}

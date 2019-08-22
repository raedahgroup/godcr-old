package nuklear

import (
	"fmt"
	"image"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type Syncer struct {
	err                error
	percentageProgress int
	report             []string
	showDetails        bool
	status             defaultsynclistener.SyncStatus
	syncError          error
}

func NewSyncer() *Syncer {
	handler := &Syncer{
		percentageProgress: 0,
		syncError:          nil,
		report: []string{
			"Starting...",
		},
		showDetails: false,
	}

	return handler
}

func (s *Syncer) startSyncing(walletMiddleware app.WalletMiddleware, masterWindow nucular.MasterWindow) {
	// begin block chain sync now so that when `Render` is called shortly after this, there'd be a report to display
	walletMiddleware.SpvSync(false, func(report *defaultsynclistener.ProgressReport) {
		progressReport := report.Read()

		s.status = progressReport.Status
		s.percentageProgress = int(progressReport.TotalSyncProgress)

		if progressReport.Status == defaultsynclistener.SyncStatusError {
			s.syncError = fmt.Errorf(progressReport.Error)
		}

		if progressReport.TotalTimeRemaining == "" {
			s.report = []string{
				fmt.Sprintf("%d%% completed.", progressReport.TotalSyncProgress),
			}
		} else {
			s.report = []string{
				fmt.Sprintf("%d%% completed, %s remaining.", progressReport.TotalSyncProgress, progressReport.TotalTimeRemaining),
			}
		}

		switch progressReport.CurrentStep {
		case defaultsynclistener.FetchingBlockHeaders:
			s.report = append(s.report, fmt.Sprintf("Fetched %d of %d block headers.",
				progressReport.FetchedHeadersCount, progressReport.TotalHeadersToFetch))
			s.report = append(s.report, fmt.Sprintf("%d%% through step 1 of 3.", progressReport.HeadersFetchProgress))

			if progressReport.DaysBehind != "" {
				s.report = append(s.report, fmt.Sprintf("Your wallet is %s behind.", progressReport.DaysBehind))
			}

		case defaultsynclistener.DiscoveringUsedAddresses:
			s.report = append(s.report, "Discovering used addresses.")
			if progressReport.AddressDiscoveryProgress > 100 {
				s.report = append(s.report, fmt.Sprintf("%d%% (over) through step 2 of 3.", progressReport.AddressDiscoveryProgress))
			} else {
				s.report = append(s.report, fmt.Sprintf("%d%% through step 2 of 3.", progressReport.AddressDiscoveryProgress))
			}

		case defaultsynclistener.ScanningBlockHeaders:
			s.report = append(s.report, fmt.Sprintf("Scanning %d of %d block headers.",
				progressReport.CurrentRescanHeight, progressReport.TotalHeadersToFetch))
			s.report = append(s.report, fmt.Sprintf("%d%% through step 3 of 3.", progressReport.HeadersFetchProgress))
		}

		// show peer count last
		if progressReport.ConnectedPeers == 1 {
			s.report = append(s.report, fmt.Sprintf("Syncing with %d peer on %s.", progressReport.ConnectedPeers, walletMiddleware.NetType()))
		} else {
			s.report = append(s.report, fmt.Sprintf("Syncing with %d peers on %s.", progressReport.ConnectedPeers, walletMiddleware.NetType()))
		}

		masterWindow.Changed()
	})
}

func (s *Syncer) isDoneSyncing() bool {
	return s.status == defaultsynclistener.SyncStatusSuccess
}

func (s *Syncer) Render(window *nucular.Window) {
	widgets.NoScrollGroupWindow("sync-page", window, func(pageWindow *widgets.Window) {
		pageWindow.Master().Style().GroupWindow.Padding = image.Point{10, 10}
		pageWindow.AddHorizontalSpace(20)
		pageWindow.AddLabelWithFont("Synchronizing", widgets.CenterAlign, styles.PageHeaderFont)

		pageWindow.PageContentWindow("sync-page-content", 10, 10, func(contentWindow *widgets.Window) {
			if s.err != nil {
				contentWindow.DisplayErrorMessage("Sync failed to start", s.err)
			} else {
				contentWindow.AddProgressBar(&s.percentageProgress, 100)

				if s.showDetails {
					for _, report := range s.report {
						contentWindow.AddLabel(report, widgets.CenterAlign)
					}
					return
				}

				contentWindow.AddLabel(s.report[0], widgets.CenterAlign)
				contentWindow.AddHorizontalSpace(20)
				contentWindow.UseFontAndResetToPrevious(styles.PageHeaderFont, func() {
					contentWindow.SelectableLabel("Tap to view information", widgets.CenterAlign, &s.showDetails)
				})
			}

			if s.syncError != nil {
				contentWindow.AddHorizontalSpace(20)
				contentWindow.DisplayErrorMessage("Sync error", s.syncError)
			}
		})
	})
}

package nuklear

import (
	"fmt"
	"image"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type Syncer struct {
	wallet             *dcrlibwallet.LibWallet
	refreshDisplay     func()
	percentageProgress int
	report             []string
	showDetails        bool
	syncError          error
}

func NewSyncer(wallet *dcrlibwallet.LibWallet, refreshDisplay func()) *Syncer {
	return &Syncer{
		wallet:             wallet,
		refreshDisplay:     refreshDisplay,
		percentageProgress: 0,
		report: []string{
			"Starting...",
		},
		showDetails: false,
		syncError:   nil,
	}
}

func (s *Syncer) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	s.refreshDisplay()
}

func (s *Syncer) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {
	s.percentageProgress = int(headersFetchProgress.TotalSyncProgress)

	s.report = []string{
		fmt.Sprintf("%d%% completed, %s remaining.", headersFetchProgress.TotalSyncProgress,
			dcrlibwallet.CalculateTotalTimeRemaining(headersFetchProgress.TotalTimeRemainingSeconds)),

		fmt.Sprintf("Fetched %d of %d block headers.", headersFetchProgress.FetchedHeadersCount,
			headersFetchProgress.TotalHeadersToFetch),

		fmt.Sprintf("%d%% through step 1 of 3.", headersFetchProgress.HeadersFetchProgress),

		fmt.Sprintf("Your wallet is %s behind.",
			dcrlibwallet.CalculateDaysBehind(headersFetchProgress.CurrentHeaderTimestamp)),
	}

	s.refreshDisplay()
}

func (s *Syncer) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {
	s.percentageProgress = int(addressDiscoveryProgress.TotalSyncProgress)

	s.report = []string{
		fmt.Sprintf("%d%% completed, %s remaining.", addressDiscoveryProgress.TotalSyncProgress,
			dcrlibwallet.CalculateTotalTimeRemaining(addressDiscoveryProgress.TotalTimeRemainingSeconds)),

		"%Discovering used addresses.",
	}

	if addressDiscoveryProgress.AddressDiscoveryProgress > 100 {
		s.report = append(s.report, fmt.Sprintf("%d%% (over) through step 2 of 3.",
			addressDiscoveryProgress.AddressDiscoveryProgress))
	} else {
		s.report = append(s.report, fmt.Sprintf("%d%% through step 2 of 3.",
			addressDiscoveryProgress.AddressDiscoveryProgress))
	}

	s.refreshDisplay()
}

func (s *Syncer) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {
	s.percentageProgress = int(headersRescanProgress.TotalSyncProgress)

	s.report = []string{
		fmt.Sprintf("%d%% completed, %s remaining.", headersRescanProgress.TotalSyncProgress,
			dcrlibwallet.CalculateTotalTimeRemaining(headersRescanProgress.TotalTimeRemainingSeconds)),

		fmt.Sprintf("Scanning %d of %d block headers.", headersRescanProgress.CurrentRescanHeight,
			headersRescanProgress.TotalHeadersToScan),

		fmt.Sprintf("%d%% through step 3 of 3.", headersRescanProgress.RescanProgress),
	}

	s.refreshDisplay()
}

func (s *Syncer) OnSyncCompleted() {
	s.percentageProgress = 100
	s.report = []string{
		"Sync completed.",
	}
	s.refreshDisplay()
}

func (s *Syncer) OnSyncCanceled(willRestart bool) {}

func (s *Syncer) OnSyncEndedWithError(err error) {
	s.syncError = err
	s.refreshDisplay()
}

func (s *Syncer) Debug(debugInfo *dcrlibwallet.DebugInfo) {}

func (s *Syncer) Render(window *nucular.Window) {
	widgets.NoScrollGroupWindow("sync-page", window, func(pageWindow *widgets.Window) {
		pageWindow.Master().Style().GroupWindow.Padding = image.Point{10, 10}
		pageWindow.AddHorizontalSpace(20)
		pageWindow.AddLabelWithFont("Synchronizing", widgets.CenterAlign, styles.PageHeaderFont)

		pageWindow.PageContentWindow("sync-page-content", 10, 10, func(contentWindow *widgets.Window) {
			contentWindow.AddProgressBar(&s.percentageProgress, 100)

			if s.syncError != nil {
				contentWindow.AddHorizontalSpace(20)
				contentWindow.DisplayErrorMessage("Sync error", s.syncError)
				return
			}

			var detailsToggleButtonText string

			if s.showDetails {
				for _, report := range s.report {
					contentWindow.AddLabel(report, widgets.CenterAlign)
				}

				// show peer count info last
				var connectedPeersInfo string
				if s.wallet.GetConnectedPeersCount() == 1 {
					connectedPeersInfo = fmt.Sprintf("Syncing with %d peer on %s.",
						s.wallet.GetConnectedPeersCount(), s.wallet.NetType())
				} else {
					connectedPeersInfo = fmt.Sprintf("Syncing with %d peer on %s.",
						s.wallet.GetConnectedPeersCount(), s.wallet.NetType())
				}
				contentWindow.AddLabel(connectedPeersInfo, widgets.CenterAlign)

				detailsToggleButtonText = "Tap to hide details"
			} else {
				contentWindow.AddLabel(s.report[0], widgets.CenterAlign)
				detailsToggleButtonText = "Tap to view details"
			}

			contentWindow.AddHorizontalSpace(20)
			contentWindow.UseFontAndResetToPrevious(styles.PageHeaderFont, func() {
				contentWindow.SelectableLabel(detailsToggleButtonText, widgets.CenterAlign, &s.showDetails)
			})
		})
	})
}

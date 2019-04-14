package pagehandlers

import (
	"fmt"
	"image"
	"github.com/raedahgroup/godcr/app/sync"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type SyncHandler struct {
	err                         error
	percentageProgress          int
	report                      []string
	showDetails          bool
	status                      sync.Status
	syncError      error
}

func (s *SyncHandler) BeforeRender(walletMiddleware app.WalletMiddleware, refreshWindowDisplay func()) {
	s.status = sync.StatusNotStarted
	s.percentageProgress = 0
	s.syncError = nil
	s.report = []string{
		"Starting...",
	}
	s.showDetails = false

	// begin block chain sync now so that when `Render` is called shortly after this, there'd be a report to display
	s.err = walletMiddleware.SyncBlockChain(false, func(syncPrivateInfo *sync.PrivateInfo) {
		syncInfo := syncPrivateInfo.Read()
		s.status = syncInfo.Status
		s.percentageProgress = int(syncInfo.TotalSyncProgress)

		if syncInfo.Status == sync.StatusError {
			s.syncError = fmt.Errorf(syncInfo.Error)
		}

		if syncInfo.TotalTimeRemaining == "" {
			s.report = []string{
				fmt.Sprintf("%d%% completed.", syncInfo.TotalSyncProgress),
			}
		} else {
			s.report = []string{
				fmt.Sprintf("%d%% completed, %s remaining.", syncInfo.TotalSyncProgress, syncInfo.TotalTimeRemaining),
			}
		}

		switch syncInfo.CurrentStep {
		case 1:
			s.report = append(s.report, fmt.Sprintf("Fetched %d of %d block headers.",
				syncInfo.FetchedHeadersCount, syncInfo.TotalHeadersToFetch))
			s.report = append(s.report, fmt.Sprintf("%d%% through step 1 of 3.", syncInfo.HeadersFetchProgress))

			if syncInfo.DaysBehind != "" {
				s.report = append(s.report, fmt.Sprintf("Your wallet is %s behind.", syncInfo.DaysBehind))
			}

		case 2:
			s.report = append(s.report, "Discovering used addresses.")
			if syncInfo.AddressDiscoveryProgress > 100 {
				s.report = append(s.report, fmt.Sprintf("%d%% (over) through step 2 of 3.", syncInfo.AddressDiscoveryProgress))
			} else {
				s.report = append(s.report, fmt.Sprintf("%d%% through step 2 of 3.", syncInfo.AddressDiscoveryProgress))
			}

		case 3:
			s.report = append(s.report, fmt.Sprintf("Scanning %d of %d block headers.",
				syncInfo.CurrentRescanHeight, syncInfo.TotalHeadersToFetch))
			s.report = append(s.report, fmt.Sprintf("%d%% through step 3 of 3.", syncInfo.HeadersFetchProgress))
		}

		// show peer count last
		if syncInfo.ConnectedPeers == 1 {
			s.report = append(s.report, fmt.Sprintf("Syncing with %d peer on %s", syncInfo.ConnectedPeers, walletMiddleware.NetType()))
		} else {
			s.report = append(s.report, fmt.Sprintf("Syncing with %d peers on %s", syncInfo.ConnectedPeers, walletMiddleware.NetType()))
		}

		refreshWindowDisplay()
	})
}

func (s *SyncHandler) Render(window *nucular.Window, changePage func(*nucular.Window, string)) {
	// change page onSyncStatusSuccess
	if s.status == sync.StatusSuccess {
		changePage(window, "overview")
		return
	}

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

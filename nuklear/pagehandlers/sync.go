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
	isRendering                 bool
	isShowingPercentageProgress bool
	percentageProgress          int
	report                      string
	status                      sync.Status
}

func (s *SyncHandler) BeforeRender() {
	s.isRendering = false
	s.report = ""
	s.isShowingPercentageProgress = false
	s.percentageProgress = 0
}

func (s *SyncHandler) Render(window *nucular.Window, wallet app.WalletMiddleware, changePage func(*nucular.Window, string)) {
	if !s.isRendering {
		s.isRendering = true
		s.syncBlockchain(window, wallet)
	}

	// change page onSyncStatusSuccess
	if s.status == sync.StatusSuccess {
		changePage(window, "overview")
		return
	}

	widgets.NoScrollGroupWindow("sync-page", window, func(pageWindow *widgets.Window) {
		pageWindow.Master().Style().GroupWindow.Padding = image.Point{10, 10}
		pageWindow.AddLabelWithFont("Synchronizing", widgets.CenterAlign, styles.PageHeaderFont)

		pageWindow.PageContentWindow("sync-page-content", 10, 10, func(contentWindow *widgets.Window) {
			if s.err != nil {
				contentWindow.DisplayErrorMessage("Error", s.err)
			} else {
				contentWindow.AddLabel(s.report, widgets.CenterAlign)
				if s.isShowingPercentageProgress {
					contentWindow.AddProgressBar(&s.percentageProgress, 100)
				}
			}
		})
	})
}

func (s *SyncHandler) syncBlockchain(window *nucular.Window, wallet app.WalletMiddleware) {
	masterWindow := window.Master()

	err := wallet.SyncBlockChainOld(&app.BlockChainSyncListener{
		SyncStarted: func() {
			s.updateStatus("Blockchain sync started...", sync.StatusInProgress)
			window.Master().Changed()
		},
		SyncEnded: func(err error) {
			if err != nil {
				s.updateStatus(fmt.Sprintf("Blockchain sync completed with error: %s", err.Error()), sync.StatusError)
			} else {
				s.updateStatus("Blockchain sync completed successfully", sync.StatusSuccess)
			}
			masterWindow.Changed()
		},
		OnHeadersFetched: func(percentageProgress int64) {
			s.updateStatusWithPercentageProgress("Blockchain sync in progress. Fetching headers (1/3)", sync.StatusInProgress, percentageProgress)
			masterWindow.Changed()
		},
		OnDiscoveredAddress: func(_ string) {
			s.updateStatus("Blockchain sync in progress. Discovering addresses (2/3)", sync.StatusInProgress)
			masterWindow.Changed()
		},
		OnRescanningBlocks: func(percentageProgress int64) {
			s.updateStatusWithPercentageProgress("Blockchain sync in progress. Rescanning blocks (3/3)", sync.StatusInProgress, percentageProgress)
			masterWindow.Changed()
		},
	}, false)

	if err != nil {
		s.err = err
		masterWindow.Changed()
	}
}

func (s *SyncHandler) updateStatusWithPercentageProgress(report string, status sync.Status, percentageProgress int64) {
	s.isShowingPercentageProgress = true
	s.report = report
	s.status = status
	s.percentageProgress = int(percentageProgress)
}

func (s *SyncHandler) updateStatus(report string, status sync.Status) {
	s.isShowingPercentageProgress = false
	s.report = report
	s.status = status
}

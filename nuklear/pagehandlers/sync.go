package pagehandlers

import (
	"fmt"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
	"image"
)

type SyncHandler struct {
	err                         error
	isRendering                 bool
	isShowingPercentageProgress bool
	percentageProgress          int
	report                      string
	status                      walletcore.SyncStatus
}

func (s *SyncHandler) BeforeRender() {
	s.isRendering = false
	s.report = ""
	s.isShowingPercentageProgress = false
	s.percentageProgress = 10
}

func (s *SyncHandler) Render(window *nucular.Window, wallet app.WalletMiddleware, changePage func(*nucular.Window, string)) {
	if !s.isRendering {
		s.isRendering = true
		s.syncBlockchain(window, wallet)
	}

	//change page onSyncStatusSuccess
	if s.status == walletcore.SyncStatusSuccess {
		changePage(window, "overview")
		return
	}

	widgets.NoScrollGroupWindow("sync-page", window, func(pageWindow *widgets.Window) {
		pageWindow.Master().Style().GroupWindow.Padding = image.Point{10, 10}
		pageWindow.AddLabelWithFont("Synchronizing", widgets.CenterAlign, styles.PageHeaderFont)

		pageWindow.PageContentWindow("sync-page-content", 10, 10, func(contentWindow *widgets.Window) {
			if s.err != nil {
				contentWindow.DisplayErrorMessage(s.err.Error())
			} else {
				contentWindow.AddLabel(s.report, widgets.CenterAlign)
				contentWindow.AddProgressBar(&s.percentageProgress, 100)

				if s.isShowingPercentageProgress {

				}
			}
		})
	})
}

func (s *SyncHandler) syncBlockchain(window *nucular.Window, wallet app.WalletMiddleware) {
	masterWindow := window.Master()

	err := wallet.SyncBlockChain(&app.BlockChainSyncListener{
		SyncStarted: func() {
			s.updateStatus("Blockchain sync started...", walletcore.SyncStatusInProgress)
			window.Master().Changed()
		},
		SyncEnded: func(err error) {
			if err != nil {
				s.updateStatus(fmt.Sprintf("Blockchain sync completed with error: %s", err.Error()), walletcore.SyncStatusError)
			} else {
				s.updateStatus("Blockchain sync completed successfully", walletcore.SyncStatusSuccess)
			}
			masterWindow.Changed()
		},
		OnHeadersFetched: func(percentageProgress int64) {
			s.updateStatusWithPercentageProgress("Blockchain sync in progress. Fetching headers (1/3)", walletcore.SyncStatusInProgress, percentageProgress)
			masterWindow.Changed()
		},
		OnDiscoveredAddress: func(_ string) {
			s.updateStatus("Blockchain sync in progress. Discovering addresses (2/3)", walletcore.SyncStatusInProgress)
			masterWindow.Changed()
		},
		OnRescanningBlocks: func(percentageProgress int64) {
			s.updateStatusWithPercentageProgress("Blockchain sync in progress. Rescanning blocks (3/3)", walletcore.SyncStatusInProgress, percentageProgress)
			masterWindow.Changed()
		},
		OnPeerConnected:    func(_ int32) {},
		OnPeerDisconnected: func(_ int32) {},
	}, false)

	if err != nil {
		s.err = err
		masterWindow.Changed()
	}
}

func (s *SyncHandler) updateStatusWithPercentageProgress(report string, status walletcore.SyncStatus, percentageProgress int64) {
	s.isShowingPercentageProgress = true
	s.report = report
	s.status = status
	s.percentageProgress = int(percentageProgress)
}

func (s *SyncHandler) updateStatus(report string, status walletcore.SyncStatus) {
	s.isShowingPercentageProgress = false
	s.report = report
	s.status = status
}

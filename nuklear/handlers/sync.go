package handlers

import (
	"fmt"
	"image/color"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type SyncHandler struct {
	err                         error
	isRendering                 bool
	isShowingPercentageProgress bool
	percentageProgress          int
	report                      string
	status                      syncStatus
}

type syncStatus uint8

const (
	syncStatusNotStarted syncStatus = iota
	syncStatusSuccess
	syncStatusError
	syncStatusInProgress
)

func (s *SyncHandler) BeforeRender() {
	s.isRendering = false
	s.report = ""
	s.isShowingPercentageProgress = false
	s.percentageProgress = 0
}

func (s *SyncHandler) Render(window *nucular.Window, wallet app.WalletMiddleware, pageChangeFunc func(string)) {
	if !s.isRendering {
		s.isRendering = true
		s.syncBlockchain(window, wallet)
	}

	// change page onSyncStatusSuccess
	if s.status == syncStatusSuccess {
		pageChangeFunc("balance")
		return
	}

	if contentWindow := helpers.NewWindow("dd", window, 0); contentWindow != nil {
		if s.err != nil {
			contentWindow.Row(50).Dynamic(1)
			contentWindow.LabelWrap(s.err.Error())
		} else {
			contentWindow.Row(40).Dynamic(1)
			contentWindow.LabelColored(s.report, "LC", color.RGBA{9, 20, 64, 255})

			if s.isShowingPercentageProgress {
				contentWindow.Row(30).Dynamic(1)
				contentWindow.Progress(&s.percentageProgress, 100, false)
			}

			contentWindow.End()
		}
	}

}

func (s *SyncHandler) syncBlockchain(window *nucular.Window, wallet app.WalletMiddleware) {
	masterWindow := window.Master()

	err := wallet.SyncBlockChain(&app.BlockChainSyncListener{
		SyncStarted: func() {
			s.updateStatus("Blockchain sync started...", syncStatusInProgress)
			window.Master().Changed()
		},
		SyncEnded: func(err error) {
			if err != nil {
				s.updateStatus(fmt.Sprintf("Blockchain sync completed with error: %s", err.Error()), syncStatusError)
			} else {
				s.updateStatus("Blockchain sync completed successfully", syncStatusSuccess)
			}
			masterWindow.Changed()
		},
		OnHeadersFetched: func(percentageProgress int64) {
			s.updateStatusWithPercentageProgress("Blockchain sync in progress. Fetching headers (1/3)", syncStatusInProgress, percentageProgress)
			masterWindow.Changed()
		},
		OnDiscoveredAddress: func(_ string) {
			s.updateStatus("Blockchain sync in progress. Discovering addresses (2/3)", syncStatusInProgress)
			masterWindow.Changed()
		},
		OnRescanningBlocks: func(percentageProgress int64) {
			s.updateStatusWithPercentageProgress("Blockchain sync in progress. Rescanning blocks (3/3)", syncStatusInProgress, percentageProgress)
			masterWindow.Changed()
		},
	}, false)

	if err != nil {
		s.err = err
		masterWindow.Changed()
	}
}

func (s *SyncHandler) updateStatusWithPercentageProgress(report string, status syncStatus, percentageProgress int64) {
	s.isShowingPercentageProgress = true
	s.report = report
	s.status = status
	s.percentageProgress = int(percentageProgress)
}

func (s *SyncHandler) updateStatus(report string, status syncStatus) {
	s.isShowingPercentageProgress = false
	s.report = report
	s.status = status
}

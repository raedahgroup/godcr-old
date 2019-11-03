package gio

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
	"github.com/raedahgroup/godcr/app"

	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
)

type Syncer struct {
	err                error
	percentageProgress int
	report             []string
	showDetails        bool
	status             defaultsynclistener.SyncStatus
	syncError          error

	informationLabel *widgets.ClickableLabel

	theme *helper.Theme
}

func NewSyncer(theme *helper.Theme) *Syncer {
	return &Syncer{
		percentageProgress: 0,
		syncError:          nil,
		report: []string{
			"Starting...",
		},
		showDetails:      false,
		theme:            theme,
		informationLabel: widgets.NewClickableLabel("Tap to view information", widgets.AlignMiddle, theme),
	}
}

func (s *Syncer) startSyncing(walletMiddleware app.WalletMiddleware, refreshWindowFunc func()) {
	// begin block chain sync now so that when `Render` is called shortly after this, there'd be a report to display
	walletMiddleware.SyncBlockChain(false, func(report *defaultsynclistener.ProgressReport) {
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

		refreshWindowFunc()
	})
}

func (s *Syncer) isDoneSyncing() bool {
	return s.status == defaultsynclistener.SyncStatusSuccess
}

func (s *Syncer) Render(ctx *layout.Context) {
	inset := layout.UniformInset(unit.Dp(10))
	inset.Layout(ctx, func() {
		if s.err != nil {
			widgets.DisplayErrorText(fmt.Sprintf("Sync failed to start: %s", s.err.Error()), s.theme, ctx)
		} else {
			widgets.NewProgressBar(&s.percentageProgress, s.theme, ctx)

			nextTopInset := float32(22)
			if s.showDetails {
				for _, report := range s.report {
					inset := layout.Inset{
						Top: unit.Dp(nextTopInset),
					}

					inset.Layout(ctx, func() {
						widgets.CenteredLabel(report, s.theme, ctx)
					})
					nextTopInset += float32(widgets.NormalLabelHeight)
					return
				}
			}

			inset := layout.Inset{
				Top: unit.Dp(nextTopInset),
			}

			inset.Layout(ctx, func() {
				clickFunc := func() {
					s.showDetails = !s.showDetails
				}
				s.informationLabel.Display(clickFunc, ctx)
			})
		}

		if s.syncError != nil {
			inset := layout.Inset{
				Top: unit.Dp(22),
			}
			inset.Layout(ctx, func() {
				widgets.DisplayErrorText(fmt.Sprintf("Sync error: %s", s.syncError.Error()), s.theme, ctx)
			})
		}
	})
}

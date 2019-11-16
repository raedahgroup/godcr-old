package gio

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
)

type Syncer struct {
	err                error
	percentageProgress int
	report             []string
	showDetails        bool
	wallet             *dcrlibwallet.LibWallet
	refreshDisplay     func()
	syncError          error
	theme              *helper.Theme

	informationLabel *widgets.ClickableLabel
}

func NewSyncer(theme *helper.Theme, wallet *dcrlibwallet.LibWallet, refreshDisplay func()) *Syncer {
	return &Syncer{
		wallet:             wallet,
		refreshDisplay:     refreshDisplay,
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

func (s *Syncer) Render(ctx *layout.Context) {
	inset := layout.UniformInset(unit.Dp(10))
	inset.Layout(ctx, func() {
		if s.err != nil {
			widgets.DisplayErrorText(fmt.Sprintf("Sync failed to start: %s", s.err.Error()), s.theme, ctx)
		} else {
			inset := layout.Inset{
				Top: unit.Dp(0),
			}
			inset.Layout(ctx, func() {
				widgets.BoldCenteredLabel("Synchronizing", s.theme, ctx)
			})

			inset = layout.Inset{
				Top: unit.Dp(22),
			}
			inset.Layout(ctx, func() {
				widgets.NewProgressBar(&s.percentageProgress, s.theme, ctx)
			})

			nextTopInset := float32(43)
			if s.showDetails {
				for _, report := range s.report {
					inset := layout.Inset{
						Top: unit.Dp(nextTopInset),
					}
					inset.Layout(ctx, func() {
						widgets.CenteredLabel(report, s.theme, ctx)
					})
					nextTopInset += float32(widgets.NormalLabelHeight)
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

				inset := layout.Inset{
					Top: unit.Dp(nextTopInset),
				}
				inset.Layout(ctx, func() {
					widgets.CenteredLabel(connectedPeersInfo, s.theme, ctx)
				})

				nextTopInset += float32(widgets.NormalLabelHeight)
				s.informationLabel.SetText("Tap to hide details")
			} else {
				s.informationLabel.SetText("Tap to view details")
			}

			inset = layout.Inset{
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

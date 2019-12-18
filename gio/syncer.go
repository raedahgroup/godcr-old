package gio

import (
	"fmt"

	"gioui.org/layout"
	//"gioui.org/unit"

	"github.com/raedahgroup/dcrlibwallet"

	//"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
)

type Syncer struct {
	err                error
	percentageProgress int
	report             []string
	showDetails        bool
	wallet             *dcrlibwallet.MultiWallet
	refreshDisplay     func()
	syncError          error
	informationLabel *widgets.ClickableLabel
}

func NewSyncer(wallet *dcrlibwallet.MultiWallet, refreshDisplay func()) *Syncer {
	return &Syncer{
		wallet:             wallet,
		refreshDisplay:     refreshDisplay,
		percentageProgress: 0,
		syncError:          nil,
		report: []string{
			"Starting...",
		},
		showDetails:      false,
		//informationLabel: widgets.NewClickableLabel("Tap to view information"),
	}
}

func (s *Syncer) OnSyncStarted() {}

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

func (s *Syncer) OnTransaction(transaction string) {}

func (s *Syncer) OnBlockAttached(walletID int, blockHeight int32) {}

func (s *Syncer) OnTransactionConfirmed(walletID int, hash string, blockHeight int32) {}

func (s *Syncer) Render(ctx *layout.Context) {
	/**inset := layout.UniformInset(unit.Dp(3))
	inset.Layout(ctx, func() {
		if s.err != nil {
			widgets.NewErrorLabel(fmt.Sprintf("Sync failed to start: %s", s.err.Error())).Draw(ctx, widgets.AlignMiddle)
		} else {
			inset := layout.Inset{
				Top: unit.Dp(5),
			}
			inset.Layout(ctx, func() {
				widgets.NewLabel("Synchronizing", 4).Draw(ctx, widgets.AlignMiddle)
			})

			inset = layout.Inset{
				Top: unit.Dp(30),
			}
			inset.Layout(ctx, func() {
				//s.widgets.ProgressBar(ctx, &s.percentageProgress)
			})

			nextTopInset := float32(47)
			if s.showDetails {
				for _, report := range s.report {
					inset := layout.Inset{
						Top: unit.Dp(nextTopInset),
					}
					inset.Layout(ctx, func() {
						widgets.NewLabel(report).Draw(ctx, widgets.AlignMiddle)
					})
					nextTopInset += float32(widgets.NormalLabelHeight)
				}

				// show peer count info last
				var connectedPeersInfo string
				if s.wallet.ConnectedPeers() == 1 {
					connectedPeersInfo = "Syncing with 1 peer"
				} else {
					connectedPeersInfo = fmt.Sprintf("Syncing with %d peers.", s.wallet.ConnectedPeers())
				}

				inset := layout.Inset{
					Top: unit.Dp(nextTopInset),
				}
				inset.Layout(ctx, func() {
					widgets.NewLabel(connectedPeersInfo).Draw(ctx, widgets.AlignMiddle)
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
				s.informationLabel.Draw(ctx, widgets.AlignMiddle, clickFunc)
			})
		}

		if s.syncError != nil {
			inset := layout.Inset{
				Top: unit.Dp(22),
			}
			inset.Layout(ctx, func() {
				widgets.NewErrorLabel(fmt.Sprintf("Sync error: %s", s.syncError.Error())).Draw(ctx, widgets.AlignMiddle)
			})
		}
	})**/
}
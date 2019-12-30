package common

import (
	"image/color"

	"strconv"
	"gioui.org/f32"
	"gioui.org/op/clip"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/text"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
)

type report struct {
	percentageProgress  float64 
	timeRemaining       string 
	daysBehind          string
}

type widgetItems struct {
	reconnectButton  *widgets.Button
	cancelButton     *widgets.Button
	progressBar      *widgets.ProgressBar
}

type Syncer struct {
	err                error
	//report             []string
	showDetails        bool
	wallet             *helper.MultiWallet
	refreshDisplay     func()
	syncError          error
	
	report           *report
	widgets          *widgetItems
}

func NewSyncer(wallet *helper.MultiWallet, refreshDisplay func()) *Syncer {
	s := &Syncer{
		wallet:             wallet,
		refreshDisplay:     refreshDisplay,
		syncError:          nil,
		showDetails:      false,
		report: &report{},
	}

	s.widgets = &widgetItems{
		progressBar: widgets.NewProgressBar(),
		reconnectButton: widgets.NewButton("Reconnect", nil).SetBorderColor(helper.GrayColor).SetBackgroundColor(helper.WhiteColor).SetColor(helper.BlackColor),
		cancelButton   : widgets.NewButton("Cancel", nil).SetBorderColor(helper.GrayColor).SetBackgroundColor(helper.WhiteColor).SetColor(helper.BlackColor),
	}



	return s
}

func (s *Syncer) OnSyncStarted() {}

func (s *Syncer) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
    s.refreshDisplay()
}

func (s *Syncer) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {
	s.report.percentageProgress = float64(headersFetchProgress.TotalSyncProgress)
	s.report.timeRemaining = dcrlibwallet.CalculateTotalTimeRemaining(headersFetchProgress.TotalTimeRemainingSeconds)
	s.report.daysBehind  = dcrlibwallet.CalculateDaysBehind(headersFetchProgress.CurrentHeaderTimestamp)


	/**s.report = []string{
		fmt.Sprintf("%d%% completed, %s remaining.", headersFetchProgress.TotalSyncProgress,
			dcrlibwallet.CalculateTotalTimeRemaining(headersFetchProgress.TotalTimeRemainingSeconds)),

		fmt.Sprintf("Fetched %d of %d block headers.", headersFetchProgress.FetchedHeadersCount,
			headersFetchProgress.TotalHeadersToFetch),

		fmt.Sprintf("%d%% through step 1 of 3.", headersFetchProgress.HeadersFetchProgress),

		fmt.Sprintf("Your wallet is %s behind.",
			dcrlibwallet.CalculateDaysBehind(headersFetchProgress.CurrentHeaderTimestamp)),
	}**/
	s.refreshDisplay()
}

func (s *Syncer) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {
	s.report.percentageProgress = float64(addressDiscoveryProgress.TotalSyncProgress)
	/**s.report = []string{
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
	}**/

	s.refreshDisplay()
}

func (s *Syncer) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {
	s.report.percentageProgress = float64(headersRescanProgress.TotalSyncProgress)
	/**s.report = []string{
		fmt.Sprintf("%d%% completed, %s remaining.", headersRescanProgress.TotalSyncProgress,
			dcrlibwallet.CalculateTotalTimeRemaining(headersRescanProgress.TotalTimeRemainingSeconds)),

		fmt.Sprintf("Scanning %d of %d block headers.", headersRescanProgress.CurrentRescanHeight,
			headersRescanProgress.TotalHeadersToScan),

		fmt.Sprintf("%d%% through step 3 of 3.", headersRescanProgress.RescanProgress),
	}**/
	s.refreshDisplay()
}

func (s *Syncer) OnSyncCompleted() {
	s.report.percentageProgress = 100
	/**s.report = []string{
		"Sync completed.",
	}**/
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
	helper.PaintArea(ctx, helper.WhiteColor, ctx.Constraints.Width.Max, 140)

	inset := layout.UniformInset(unit.Dp(15))
	inset.Layout(ctx, func(){
		layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
			layout.Rigid(func(){
				widgets.NewLabel("Wallet Status").
					SetColor(helper.GrayColor).
					SetSize(4).
					SetWeight(text.Bold).
					Draw(ctx)
			}),
			layout.Flexed(1, func(){
				layout.Align(layout.NE).Layout(ctx, func(){
					s.drawWalletStatus(ctx)
				})
			}),
		)


		inset := layout.Inset{
			Top: unit.Dp(25),
		}

		inset.Layout(ctx, func(){
			if !s.wallet.IsSynced() {
				s.drawIsSyncingStatus(ctx)
				//s.drawNotSyncedStatus(ctx)
			} else {
				s.drawIsSyncingStatus(ctx)
			}
		})
	})

	


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


func (s *Syncer) drawWalletStatus(ctx *layout.Context) {
	var indicatorColor color.RGBA 
	var statusText string 

	if s.IsOnline() {
		indicatorColor = helper.DecredGreenColor 
		statusText = "Online"
	} else {
		indicatorColor = helper.DecredOrangeColor 
		statusText = "Offline"
	}

	indicatorSize := float32(8)
	radius := indicatorSize * .5 

	layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
		layout.Rigid(func(){
			inset := layout.Inset{
				Top: unit.Dp(3.5),
				Right: unit.Dp(5),
			}
			inset.Layout(ctx, func(){
				clip.Rect{
					Rect: f32.Rectangle{
						Max: f32.Point{
							X: indicatorSize,
							Y: indicatorSize,
						},
					},
					NE: radius,
					NW: radius,
					SE: radius,
					SW: radius,
				}.Op(ctx.Ops).Add(ctx.Ops)
				helper.Fill(ctx, indicatorColor, int(indicatorSize), int(indicatorSize))
			})
		}),
		layout.Rigid(func(){
			widgets.NewLabel(statusText).
				SetSize(4).
				SetColor(helper.GrayColor).
				SetWeight(text.Bold).
				Draw(ctx)
		}),
	)
}


func (s *Syncer) IsOnline() bool {
	if !s.wallet.IsSynced() || !s.wallet.IsSyncing() {
		return false
	}
	return true
}

func (s *Syncer) drawNotSyncedStatus(ctx *layout.Context) {
	layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
		layout.Rigid(func(){
			widgets.CancelIcon.SetColor(helper.DangerColor).Draw(ctx, 25)
		}),
		layout.Rigid(func(){
			inset := layout.Inset{
				Left: unit.Dp(40),
			}
			inset.Layout(ctx, func(){
				widgets.NewLabel("Not Synced").
					SetSize(5).
					SetColor(helper.BlackColor).
					SetWeight(text.Bold).
					Draw(ctx)
			})
		}),
		layout.Flexed(1, func(){
			layout.Align(layout.NE).Layout(ctx, func(){
				ctx.Constraints.Height.Max = 35
				ctx.Constraints.Width.Max = 130

				s.widgets.reconnectButton.Draw(ctx, func(){

				})
			})
		}),
	)

	inset := layout.Inset{
		Top: unit.Dp(35),
		Left: unit.Dp(40),
	}
	inset.Layout(ctx, func(){
		lowestBlock := s.wallet.GetLowestBlock()
		lowestBlockHeight := int32(-1)
		
		if lowestBlock != nil {
			lowestBlockHeight = lowestBlock.Height
		}

		txt := "Synced to block " + strconv.Itoa(int(lowestBlockHeight)) + " - " + s.report.daysBehind
		widgets.NewLabel(txt).
			SetSize(4).
			SetColor(helper.GrayColor).
			Draw(ctx)
	})

	inset = layout.Inset{
		Top: unit.Dp(65),
		Left: unit.Dp(40),
	}
	inset.Layout(ctx, func(){
		widgets.NewLabel("No connected peers").
			SetSize(4).
			SetColor(helper.GrayColor).
			Draw(ctx)
	})
}

func (s *Syncer) drawIsSyncingStatus(ctx *layout.Context) {
	layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
		layout.Rigid(func(){
			inset := layout.Inset{
				Top: unit.Dp(3),
			}
			inset.Layout(ctx, func(){
				helper.PaintCircle(ctx, helper.DecredGreenColor, 15)
			})
		}),
		layout.Rigid(func(){
			inset := layout.Inset{
				Left: unit.Dp(30),
			}
			inset.Layout(ctx, func(){
				widgets.NewLabel("Syncing...").
					SetSize(5).
					SetColor(helper.BlackColor).
					SetWeight(text.Bold).
					Draw(ctx)
			})
		}),
		layout.Flexed(1, func(){
			layout.Align(layout.NE).Layout(ctx, func(){
				ctx.Constraints.Height.Max = 35
				ctx.Constraints.Width.Max = 130

				s.widgets.cancelButton.Draw(ctx, func(){

				})
			})
		}),
	)

	inset := layout.Inset{
		Top: unit.Dp(38),
	}
	inset.Layout(ctx, func(){
		s.widgets.progressBar.
			SetHeight(10).
			SetBackgroundColor(helper.GrayColor).
			Draw(ctx, &s.report.percentageProgress)
	})

	inset = layout.Inset{
		Top: unit.Dp(50),
	}
	inset.Layout(ctx, func(){
		layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
			layout.Rigid(func(){
				txt := strconv.Itoa(int(s.report.percentageProgress)) + "%"
				widgets.NewLabel(txt).SetColor(helper.BlackColor).SetSize(4).SetWeight(text.Bold).Draw(ctx)
			}),
			layout.Flexed(1, func(){
				layout.Align(layout.NE).Layout(ctx, func(){
					widgets.NewLabel(s.report.timeRemaining).
						SetColor(helper.BlackColor).
						SetSize(4).
						SetWeight(text.Bold).
						Draw(ctx)
				})
			}),
		)
	})
}


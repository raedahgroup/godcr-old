package pages

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

var SyncDone bool
var Wallet godcrApp.WalletMiddleware

func ShowSyncWindow(wallet godcrApp.WalletMiddleware, window fyne.Window, App fyne.App) fyne.CanvasObject {
	Wallet = wallet
	go func() {
		//wait for sync to be done before calling overview
		for SyncDone == false {
			sec, _ := time.ParseDuration("1s")
			time.Sleep(sec)
		}
		window.SetContent(Menu(OverviewPage(window, App), window, App))
	}()
	go func() {
		for {
			//wait for 5 seconds to get peer info and block height
			sec, _ := time.ParseDuration("5s")
			time.Sleep(sec)
			info, _ := wallet.WalletConnectionInfo()
			if info.PeersConnected <= 1 {
				PeerConn.Text = strconv.Itoa(int(info.PeersConnected)) + " Peer Connected"
			} else {
				PeerConn.Text = strconv.Itoa(int(info.PeersConnected)) + " Peers Connected"
			}
			BlkHeight.Text = strconv.Itoa(int(info.LatestBlock)) + " Blocks Connected"
			canvas.Refresh(PeerConn)
			canvas.Refresh(BlkHeight)
		}
	}()
	progressBar := widget.NewProgressBar()
	progressBar.Min = 0
	progressBar.Max = 100

	var fullSyncReport string

	reportLabel := widget.NewLabel("")
	widget.Refresh(reportLabel)
	reportLabel.Hide()
	reportLabel.Alignment = fyne.TextAlignCenter
	var infoButton *widget.Button

	infoButton = widget.NewButton("Tap to view informations", func() {
		if infoButton.Text == "Tap to view informations" {

			widget.Refresh(reportLabel)
			reportLabel.Show()
			infoButton.SetText("Tap to hide informations")

		} else {

			widget.Refresh(reportLabel)
			reportLabel.Hide()
			infoButton.SetText("Tap to view informations")
			
		}
	})

	wallet.SyncBlockChain(false, func(report *defaultsynclistener.ProgressReport) {
		progressReport := report.Read()

		progressBar.SetValue(float64(progressReport.TotalSyncProgress))

		if progressReport.Status == defaultsynclistener.SyncStatusSuccess {
			SyncDone = true
			return
		}

		stringReport := strings.Builder{}
		if progressReport.TotalTimeRemaining == "" {
			stringReport.WriteString(fmt.Sprintf("%d%% completed.\n", progressReport.TotalSyncProgress))
		} else {
			stringReport.WriteString(fmt.Sprintf("%d%% completed, %s remaining.\n", progressReport.TotalSyncProgress, progressReport.TotalTimeRemaining))
		}

		reportLabel.SetText(strings.TrimSpace(stringReport.String()))
		widget.Refresh(reportLabel)

		switch progressReport.CurrentStep {
		case defaultsynclistener.FetchingBlockHeaders:
			stringReport.WriteString(fmt.Sprintf("Fetched %d of %d block headers.\n", progressReport.FetchedHeadersCount, progressReport.TotalHeadersToFetch))
			stringReport.WriteString(fmt.Sprintf("%d%% through step 1 of 3.\n", progressReport.HeadersFetchProgress))

			if progressReport.DaysBehind != "" {
				stringReport.WriteString(fmt.Sprintf("Your wallet is %s behind.\n", progressReport.DaysBehind))
			}
		case defaultsynclistener.DiscoveringUsedAddresses:
			stringReport.WriteString("Discovering used addresses.\n")
			if progressReport.AddressDiscoveryProgress > 100 {
				stringReport.WriteString(fmt.Sprintf("%d%% (over) through step 2 of 3.\n", progressReport.AddressDiscoveryProgress))
			} else {
				stringReport.WriteString(fmt.Sprintf("%d%% through step 2 of 3.\n", progressReport.AddressDiscoveryProgress))
			}

		case defaultsynclistener.ScanningBlockHeaders:
			stringReport.WriteString(fmt.Sprintf("Scanning %d of %d block headers.\n", progressReport.CurrentRescanHeight,
				progressReport.TotalHeadersToFetch))
			stringReport.WriteString(fmt.Sprintf("%d%% through step 3 of 3.\n", progressReport.HeadersFetchProgress))
		}

		// show peer count last
		netType := wallet.NetType()
		if progressReport.ConnectedPeers == 1 {
			stringReport.WriteString(fmt.Sprintf("Syncing with %d peer on %s.\n", progressReport.ConnectedPeers, netType))
		} else {
			stringReport.WriteString(fmt.Sprintf("Syncing with %d peers on %s.\n", progressReport.ConnectedPeers, netType))
		}

		fullSyncReport = stringReport.String()
		reportLabel.SetText(fullSyncReport)
		widget.Refresh(reportLabel)
	})

	return widget.NewVBox(
		widgets.NewVSpacer(10),
		progressBar,
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(infoButton.MinSize()), infoButton),
		reportLabel)
}

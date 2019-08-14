package pages

import (
	"context"
	"strconv"
	"time"

	"fyne.io/fyne/theme"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
)

type menuPageData struct {
	peerConn  *widget.Label
	blkHeight *widget.Label
	//there might be situations we would want to get the particular opened tab
	tabs *widget.TabContainer
	//when theme changes, this updates the canvas text
	alphaTheme uint8
}

var menu menuPageData

//this updates peerconn and blkheight
func statusUpdates(wallet godcrApp.WalletMiddleware) {
	info, _ := wallet.WalletConnectionInfo()

	if info.PeersConnected <= 1 {
		menu.peerConn.SetText(strconv.Itoa(int(info.PeersConnected)) + " Peer Connected")
	} else {
		menu.peerConn.SetText(strconv.Itoa(int(info.PeersConnected)) + " Peers Connected")
	}

	menu.blkHeight.SetText(strconv.Itoa(int(info.LatestBlock)) + " Blocks Connected")
}

func pageNotImplemented() fyne.CanvasObject {
	label := widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	return label
}

func menuPage(ctx context.Context, wallet godcrApp.WalletMiddleware, appSettings *config.Settings, fyneApp fyne.App, window fyne.Window) fyne.CanvasObject {
	if fyneApp.Settings().Theme() == theme.LightTheme() {
		menu.alphaTheme = 255
	} else {
		menu.alphaTheme = 0
	}

	menu.peerConn = widget.NewLabel("")
	menu.blkHeight = widget.NewLabel("")

	menu.tabs = widget.NewTabContainer(
		widget.NewTabItem("Overview", overviewPage(wallet)),
		widget.NewTabItem("History", pageNotImplemented()),
		widget.NewTabItem("Send", pageNotImplemented()),
		widget.NewTabItem("Receive", receivePage(wallet, window)),
		widget.NewTabItem("Staking", pageNotImplemented()),
		widget.NewTabItem("Accounts", accountPage(wallet, appSettings, window)),
		widget.NewTabItem("Security", pageNotImplemented()),
		widget.NewTabItem("Settings", settingsPage(fyneApp)),
		widget.NewTabItem("Exit", exit(ctx, fyneApp, window)))
	menu.tabs.SetTabLocation(widget.TabLocationLeading)

	//this would update all labels for all pages every seconds, all objects to be updated should be placed here
	go func() {
		for {
			//update only when the user is on the page
			if menu.tabs.CurrentTabIndex() == 0 {
				overviewPageUpdates(wallet)
			} else if menu.tabs.CurrentTabIndex() == 3 {
				receivePageUpdates(wallet)
			}
			statusUpdates(wallet)
			time.Sleep(time.Second * 1)
		}
	}()

	//where peerConn and blkHeight are the realtime status texts
	status := widget.NewVBox(menu.peerConn, menu.blkHeight)
	data := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, status, menu.tabs, nil), menu.tabs, status)

	return data
}

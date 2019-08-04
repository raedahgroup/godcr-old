package pages

import (
	"context"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/theme"

	"fyne.io/fyne/canvas"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
)

type menuPageData struct {
	alphaTheme uint8
	peerConn   *canvas.Text
	blkHeight  *canvas.Text
	//there might be situations we would want to get the particular opened tab
	tabs *widget.TabContainer
}

var menu menuPageData

//this updates peerconn and blkheight
func statusUpdates(wallet godcrApp.WalletMiddleware) {
	info, _ := wallet.WalletConnectionInfo()

	if info.PeersConnected == 1 {
		menu.peerConn.Text = (strconv.Itoa(int(info.PeersConnected)) + " Peer Connected")
		menu.peerConn.Color = color.RGBA{11, 156, 49, menu.alphaTheme}
	} else if info.PeersConnected > 1 {
		menu.peerConn.Text = (strconv.Itoa(int(info.PeersConnected)) + " Peers Connected")
		menu.peerConn.Color = color.RGBA{11, 156, 49, menu.alphaTheme}
	} else {
		menu.peerConn.Text = ("No Peer Connected")
		menu.peerConn.Color = color.RGBA{255, 0, 0, menu.alphaTheme}
	}
	canvas.Refresh(menu.peerConn)

	menu.blkHeight.Text = (strconv.Itoa(int(info.LatestBlock)) + " Blocks Connected")
	menu.blkHeight.Color = color.RGBA{11, 156, 49, menu.alphaTheme}
	canvas.Refresh(menu.blkHeight)
}

func pageNotImplemented() fyne.CanvasObject {
	label := widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	return label
}

func menuPage(ctx context.Context, wallet godcrApp.WalletMiddleware, app fyne.App, window fyne.Window) fyne.CanvasObject {
	if app.Settings().Theme() == theme.LightTheme() {
		menu.alphaTheme = 255
	} else {
		menu.alphaTheme = 200
	}
	menu.peerConn = canvas.NewText("", color.RGBA{11, 156, 49, menu.alphaTheme})
	menu.peerConn.TextStyle = fyne.TextStyle{Bold: true}
	menu.blkHeight = canvas.NewText("", color.RGBA{11, 156, 49, menu.alphaTheme})
	menu.blkHeight.TextStyle = fyne.TextStyle{Bold: true}

	menu.tabs = widget.NewTabContainer(
		widget.NewTabItem("Overview", overviewPage(wallet)),
		widget.NewTabItem("History", pageNotImplemented()),
		widget.NewTabItem("Send", sendPage(wallet, window)), //send(window)),
		widget.NewTabItem("Receive", receivePage(wallet, window)),
		widget.NewTabItem("Staking", pageNotImplemented()),
		widget.NewTabItem("Accounts", pageNotImplemented()),
		widget.NewTabItem("Security", pageNotImplemented()),
		widget.NewTabItem("Settings", settingsPage(app)),
		widget.NewTabItem("Exit", exit(ctx, app, window)))
	menu.tabs.SetTabLocation(widget.TabLocationLeading)

	//this would update all labels for all pages every seconds, all objects to be updated should be placed here
	go func() {
		for {
			//update only when the user is on the page
			if menu.tabs.CurrentTabIndex() == 0 {
				overviewUpdates(wallet)
			} else if menu.tabs.CurrentTabIndex() == 2 {
				sendPageUpdates(wallet)
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

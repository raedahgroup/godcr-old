package pages

import (
	"context"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
)

type menuPageData struct {
	peerConn  *widget.Label
	blkHeight *widget.Label
	//there might be situations we would want to get the particular opened tab
	tabs *widget.TabContainer
}

var menu menuPageData

func pageNotImplemented() fyne.CanvasObject {
	label := widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	return label
}

func menuPage(ctx context.Context, wallet godcrApp.WalletMiddleware, app fyne.App, window fyne.Window) fyne.CanvasObject {
	menu.peerConn = widget.NewLabel("")
	menu.blkHeight = widget.NewLabel("")

	menu.tabs = widget.NewTabContainer(
		widget.NewTabItem("Overview", overviewPage(wallet)),
		widget.NewTabItem("History", pageNotImplemented()),
		widget.NewTabItem("Send", pageNotImplemented()),
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
			} else if menu.tabs.CurrentTabIndex() == 3 {
				receiveUpdates(wallet)
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

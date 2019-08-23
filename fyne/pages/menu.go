package pages

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
)

type menuPageData struct {
	peerConn  *widget.Label
	blkHeight *widget.Label
	// there might be situations we would want to get the particular opened tab
	tabs *widget.TabContainer
}

var menu menuPageData

func pageNotImplemented() fyne.CanvasObject {
	label := widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	return label
}

func menuPage(ctx context.Context, wallet godcrApp.WalletMiddleware, fyneApp fyne.App, window fyne.Window) fyne.CanvasObject {
	overviewFile, err := ioutil.ReadFile("./fyne/pages/png/overview.png")
	if err != nil {
		log.Fatalln("overview png file missing", err)
	}
	historyFile, err := ioutil.ReadFile("./fyne/pages/png/history.png")
	if err != nil {
		log.Fatalln("history png file missing", err)
	}
	sendFile, err := ioutil.ReadFile("./fyne/pages/png/send.png")
	if err != nil {
		log.Fatalln("send png file missing", err)
	}
	receiveFile, err := ioutil.ReadFile("./fyne/pages/png/receive.png")
	if err != nil {
		log.Fatalln("receive png file missing", err)
	}
	accountsFile, err := ioutil.ReadFile("./fyne/pages/png/account.png")
	if err != nil {
		log.Fatalln("account png file missing", err)
	}
	moreFile, err := ioutil.ReadFile("./fyne/pages/png/more.png")
	if err != nil {
		log.Fatalln("security png file missing", err)
	}
	exitFile, err := ioutil.ReadFile("./fyne/pages/png/exit.png")
	if err != nil {
		log.Fatalln("exit png file missing", err)
	}
	stakingFile, err := ioutil.ReadFile("./fyne/pages/png/stake.png")
	if err != nil {
		log.Fatalln("staking png file missing", err)
	}

	menu.peerConn = widget.NewLabel("")
	menu.blkHeight = widget.NewLabel("")

	menu.tabs = widget.NewTabContainer(
		widget.NewTabItemWithIcon("Overview", fyne.NewStaticResource("Overview", overviewFile), overviewPage(wallet, fyneApp)),
		widget.NewTabItemWithIcon("History", fyne.NewStaticResource("History", historyFile), pageNotImplemented()),
		widget.NewTabItemWithIcon("Send", fyne.NewStaticResource("Send", sendFile), pageNotImplemented()),
		widget.NewTabItemWithIcon("Receive", fyne.NewStaticResource("Receive", receiveFile), receivePage(wallet, window)),
		widget.NewTabItemWithIcon("Accounts", fyne.NewStaticResource("Accounts", accountsFile), pageNotImplemented()),
		widget.NewTabItemWithIcon("Staking", fyne.NewStaticResource("Staking", stakingFile), stakingPage(wallet)),
		widget.NewTabItemWithIcon("Tickets", fyne.NewStaticResource("Tickets", stakingFile), initTicketPage()),
		widget.NewTabItemWithIcon("More", fyne.NewStaticResource("More", moreFile), morePage(wallet, fyneApp)),
		widget.NewTabItemWithIcon("Exit", fyne.NewStaticResource("Exit", exitFile), exit(ctx, fyneApp, window)))
	menu.tabs.SetTabLocation(widget.TabLocationLeading)

	// would update all labels for all pages every seconds, all objects to be updated should be placed here
	go func() {
		var lastActiveTabIndex = 0
		for {
			// update only when the user returns to the page after leaving previously
			switch menu.tabs.CurrentTabIndex() {
			case lastActiveTabIndex:
				// still on the same tab, no need to refresh
				break
			case 0:
				overviewPageUpdates(wallet)
			case 3:
				receivePageUpdates(wallet)
			case 5:
				stakingPageReloadData(wallet)
			case 6:
				ticketPageUpdates(wallet)
			}

			statusUpdates(wallet)
			lastActiveTabIndex = menu.tabs.CurrentTabIndex()
			time.Sleep(time.Second * 1)
		}
	}()

	// where peerConn and blkHeight are the realtime status texts
	status := widget.NewVBox(menu.peerConn, menu.blkHeight)
	data := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, status, menu.tabs, nil), menu.tabs, status)

	return data
}

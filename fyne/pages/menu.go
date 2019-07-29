package pages

import (
	"io/ioutil"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
)

var (
	peerConn  *widget.Label
	blkHeight *widget.Label
	//there might be situations we would want to get the particular opened tab
	tabs *widget.TabContainer
)

func init() {
	peerConn = widget.NewLabel("No peers connected")
	blkHeight = widget.NewLabel("No blocks fetched")

}

func pageNotImplemented() fyne.CanvasObject {
	label := widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	return label
}

func Menu(wallet godcrApp.WalletMiddleware, app fyne.App, window fyne.Window) fyne.CanvasObject {
	overviewFile, err := ioutil.ReadFile("./fyne/pages/png/overview.png")
	if err != nil {
		log.Fatalln("overview file missing", err)
	}
	historyFile, err := ioutil.ReadFile("./fyne/pages/png/block.png")
	if err != nil {
		log.Fatalln("history file missing", err)
	}
	sendFile, err := ioutil.ReadFile("./fyne/pages/png/send.png")
	if err != nil {
		log.Fatalln("send file missing", err)
	}
	receiveFile, err := ioutil.ReadFile("./fyne/pages/png/receive.png")
	if err != nil {
		log.Fatalln("receive file missing", err)
	}
	stakingFile, err := ioutil.ReadFile("./fyne/pages/png/stakeyBaby.png")
	if err != nil {
		log.Fatalln("staking file missing", err)
	}
	accountsFile, err := ioutil.ReadFile("./fyne/pages/png/wallet.png")
	if err != nil {
		log.Fatalln("accounts file missing", err)
	}
	securityFile, err := ioutil.ReadFile("./fyne/pages/png/info.png")
	if err != nil {
		log.Fatalln("security file missing", err)
	}
	settingsFile, err := ioutil.ReadFile("./fyne/pages/png/gears.png")
	if err != nil {
		log.Fatalln("settings file missing", err)
	}

	tabs = widget.NewTabContainer(
		widget.NewTabItemWithIcon("Overview", fyne.NewStaticResource("Overview", overviewFile), overviewPage(wallet)),
		widget.NewTabItemWithIcon("History", fyne.NewStaticResource("History", historyFile), pageNotImplemented()),
		widget.NewTabItemWithIcon("Send", fyne.NewStaticResource("Send", sendFile), pageNotImplemented()),
		widget.NewTabItemWithIcon("Receive", fyne.NewStaticResource("Receive", receiveFile), pageNotImplemented()),
		widget.NewTabItemWithIcon("Staking", fyne.NewStaticResource("Staking", stakingFile), pageNotImplemented()),
		widget.NewTabItemWithIcon("Accounts", fyne.NewStaticResource("Accounts", accountsFile), pageNotImplemented()),
		widget.NewTabItemWithIcon("Security", fyne.NewStaticResource("Security", securityFile), pageNotImplemented()),
		widget.NewTabItemWithIcon("Settings", fyne.NewStaticResource("Settings", settingsFile), settingsPage(app)))
	tabs.SetTabLocation(widget.TabLocationLeading)

	//where peerConn and blkHeight are the realtime status texts
	status := widget.NewVBox(peerConn, blkHeight)
	menu := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, status, tabs, nil), tabs, status)
	return menu
}

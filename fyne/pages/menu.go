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
		log.Fatalln("overviewFile file missing", err)
	}
	historyFile, err := ioutil.ReadFile("./fyne/pages/png/block.png")
	if err != nil {
		log.Fatalln("historyFile file missing", err)
	}
	sendFile, err := ioutil.ReadFile("./fyne/pages/png/send.png")
	if err != nil {
		log.Fatalln("sendFile file missing", err)
	}
	receiveFile, err := ioutil.ReadFile("./fyne/pages/png/receive.png")
	if err != nil {
		log.Fatalln("receiveFile file missing", err)
	}
	stakingFile, err := ioutil.ReadFile("./fyne/pages/png/stakeyBaby.png")
	if err != nil {
		log.Fatalln("stakingFile file missing", err)
	}
	accountsFile, err := ioutil.ReadFile("./fyne/pages/png/wallet.png")
	if err != nil {
		log.Fatalln("accountsFile file missing", err)
	}
	securityFile, err := ioutil.ReadFile("./fyne/pages/png/info.png")
	if err != nil {
		log.Fatalln("securityFile file missing", err)
	}
	settingsFile, err := ioutil.ReadFile("./fyne/pages/png/gears.png")
	if err != nil {
		log.Fatalln("settingsFile file missing", err)
	}

	tabs = widget.NewTabContainer(
		widget.NewTabItemWithIcon("Overview", fyne.NewStaticResource("Overview", overviewFile), overviewPage(wallet)),
		widget.NewTabItemWithIcon("History", fyne.NewStaticResource("History", historyFile), widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})),
		widget.NewTabItemWithIcon("Send", fyne.NewStaticResource("Send", sendFile), widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})),
		widget.NewTabItemWithIcon("Receive", fyne.NewStaticResource("Receive", receiveFile), widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})),
		widget.NewTabItemWithIcon("Staking", fyne.NewStaticResource("Staking", stakingFile), widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})),
		widget.NewTabItemWithIcon("Accounts", fyne.NewStaticResource("Accounts", accountsFile), widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})),
		widget.NewTabItemWithIcon("Security", fyne.NewStaticResource("Security", securityFile), widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})),
		widget.NewTabItemWithIcon("Settings", fyne.NewStaticResource("Settings", settingsFile), settingsPage(app)))
	tabs.SetTabLocation(widget.TabLocationLeading)

	//where peerConn and blkHeight are the realtime status texts
	status := widget.NewVBox(peerConn, blkHeight)
	menu := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, status, tabs, nil), tabs, status)
	return menu
}

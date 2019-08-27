package pages

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/raedahgroup/godcr/fyne/widgets"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
)

type menuPageData struct {
	peerConn  *widget.Label
	blkHeight *widget.Label
	// there might be situations we would want to get the particular opened tab
	tabs *widget.TabContainer
	//when theme changes, this updates the canvas text
	alphaTheme uint8
}

type pageContainer struct {
	container *widget.Box
}

var menu menuPageData

func resetPages(exempt int, window fyne.Window) {
	for i := 0; i < len(menu.tabs.Items); i++ {
		if i == exempt {
			continue
		}
		a, ok := interface{}(menu.tabs.Items[i].Content).(*widget.Box)
		if !ok {
			continue
		}

		a.Children = widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("fetching data...", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Bold: true, Italic: true})).Children
		widget.Refresh(a)
	}
}

func pageNotImplemented() fyne.CanvasObject {
	label := widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	return widget.NewHBox(widgets.NewHSpacer(10), label)
}

func menuPage(ctx context.Context, wallet godcrApp.WalletMiddleware, appSettings *config.Settings, fyneApp fyne.App, window fyne.Window) fyne.CanvasObject {
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

	if fyneApp.Settings().Theme() == theme.LightTheme() {
		menu.alphaTheme = 255
	} else {
		menu.alphaTheme = 0
	}

	overviewPageContainer.container = widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("fetching data...", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Bold: true, Italic: true}))
	historyPageContainer.container = widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("fetching data...", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Bold: true, Italic: true}))
	receivePageContainer.container = widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("fetching data...", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Bold: true, Italic: true}))
	stakingPageContainer.container = widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("fetching data...", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Bold: true, Italic: true}))
	accountPageContainer.container = widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("fetching data...", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Bold: true, Italic: true}))

	overviewPage(wallet, fyneApp)
	menu.tabs = widget.NewTabContainer(
		widget.NewTabItemWithIcon("Overview", fyne.NewStaticResource("Overview", overviewFile), overviewPageContainer.container), //overviewPage(wallet, fyneApp)),
		widget.NewTabItemWithIcon("History", fyne.NewStaticResource("History", historyFile), historyPageContainer.container),
		widget.NewTabItemWithIcon("Send", fyne.NewStaticResource("Send", sendFile), pageNotImplemented()),
		widget.NewTabItemWithIcon("Receive", fyne.NewStaticResource("Receive", receiveFile), receivePageContainer.container),    //receivePage(wallet, window)),
		widget.NewTabItemWithIcon("Accounts", fyne.NewStaticResource("Accounts", accountsFile), accountPageContainer.container), // accountPage(wallet, appSettings, window)),
		widget.NewTabItemWithIcon("Staking", fyne.NewStaticResource("Staking", stakingFile), stakingPageContainer.container),    //stakingPage(wallet)),
		widget.NewTabItemWithIcon("More", fyne.NewStaticResource("More", moreFile), morePage(wallet, fyneApp)),
		widget.NewTabItemWithIcon("Exit", fyne.NewStaticResource("Exit", exitFile), exit(ctx, fyneApp, window)))
	menu.tabs.SetTabLocation(widget.TabLocationLeading)

	// would update all labels for all pages every seconds, all objects to be updated should be placed here
	go func() {
		var currentPage = 0

		for {
			fmt.Println(window.Content())
			// load contents to page when user is on the page
			// update only when the user is on the page
			if menu.tabs.CurrentTabIndex() == 0 {
				if currentPage != 0 {
					history = historyPageData{}
					overviewPage(wallet, fyneApp)
					resetPages(0, window)
					currentPage = 0

				}
				//overviewPageUpdates(wallet)
			} else if menu.tabs.CurrentTabIndex() == 1 {
				if currentPage != 1 {
					overview = overviewPageData{}
					historyPage(wallet, window)
					resetPages(1, window)
					currentPage = 1
				}
				historyPageUpdates(wallet, window)
			} else if menu.tabs.CurrentTabIndex() == 2 {
				if currentPage != 2 {
					resetPages(2, window)
					currentPage = 2
				}
			} else if menu.tabs.CurrentTabIndex() == 3 {
				if currentPage != 3 {
					overview = overviewPageData{}
					history = historyPageData{}
					receivePage(wallet, window)
					resetPages(3, window)
					currentPage = 3
				}
			} else if menu.tabs.CurrentTabIndex() == 4 {
				if currentPage != 4 {
					overview = overviewPageData{}
					history = historyPageData{}
					accountPage(wallet, appSettings, window)
					resetPages(4, window)
					currentPage = 4
				}
			} else if menu.tabs.CurrentTabIndex() == 5 {
				if currentPage != 5 {
					overview = overviewPageData{}
					history = historyPageData{}
					stakingPage(wallet)
					resetPages(5, window)
					currentPage = 5
				}
			} else if menu.tabs.CurrentTabIndex() == 6 || menu.tabs.CurrentTabIndex() == 7 {
				overview = overviewPageData{}
				history = historyPageData{}
				resetPages(6, window)
				currentPage = 6
			}
			statusUpdates(wallet)

			time.Sleep(time.Second * 1)
		}
	}()

	// where peerConn and blkHeight are the realtime status texts
	status := widget.NewVBox(menu.peerConn, menu.blkHeight)
	data := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, status, menu.tabs, nil), menu.tabs, status)

	return data
}

package pages

import (
	"context"
	"fmt"
	"strconv"
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

		a.Children = widget.NewHBox().Children
		widget.Refresh(a)
	}
}

func pageNotImplemented() fyne.CanvasObject {
	label := widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	return widget.NewHBox(widgets.NewHSpacer(10), label)
}

func statusUpdates(wallet godcrApp.WalletMiddleware) {
	info, err := wallet.WalletConnectionInfo()
	if err != nil {
		widget.Refresh(overview.errorLabel)
		overview.errorLabel.SetText(err.Error())
	}

	if info.PeersConnected <= 1 {
		menu.peerConn.SetText(strconv.Itoa(int(info.PeersConnected)) + " Peer Connected")
	} else {
		menu.peerConn.SetText(strconv.Itoa(int(info.PeersConnected)) + " Peers Connected")
	}

	menu.blkHeight.SetText(strconv.Itoa(int(info.LatestBlock)) + " Blocks Connected")
}

func menuPage(ctx context.Context, wallet godcrApp.WalletMiddleware, appSettings *config.Settings, fyneApp fyne.App, window fyne.Window) fyne.CanvasObject {
	menu.peerConn = widget.NewLabel("")
	menu.blkHeight = widget.NewLabel("")

	if fyneApp.Settings().Theme() == theme.LightTheme() {
		menu.alphaTheme = 255
	} else {
		menu.alphaTheme = 0
	}

	overviewPageContainer.container = widget.NewHBox()
	historyPageContainer.container = widget.NewHBox()
	receivePageContainer.container = widget.NewHBox()
	stakingPageContainer.container = widget.NewHBox()
	accountPageContainer.container = widget.NewHBox()

	overviewPage(wallet, fyneApp)
	menu.tabs = widget.NewTabContainer(
		widget.NewTabItemWithIcon("Overview", overviewIcon, overviewPageContainer.container),
		widget.NewTabItemWithIcon("History", historyIcon, historyPageContainer.container),
		widget.NewTabItemWithIcon("Send", sendIcon, pageNotImplemented()),
		widget.NewTabItemWithIcon("Receive", receiveIcon, receivePageContainer.container),
		widget.NewTabItemWithIcon("Accounts", accountIcon, accountPageContainer.container),
		widget.NewTabItemWithIcon("Staking", stakingIcon, stakingPageContainer.container),
		widget.NewTabItemWithIcon("More", moreIcon, morePage(wallet, fyneApp)),
		widget.NewTabItemWithIcon("Exit", exitIcon, exit(ctx, fyneApp, window)))
	menu.tabs.SetTabLocation(widget.TabLocationLeading)

	// This goroutine tracks tab changes, and deallocates unneeded tab memory contents.
	go func() {
		// Todo: notImplemented function should be removed when all page has been implemented
		notImplemented := func(page int) {
			if (page + 1) > len(menu.tabs.Items) {
				fmt.Println("page not available in tab")
				return
			}
			a, ok := interface{}(menu.tabs.Items[page].Content).(*widget.Box)
			if !ok {
				return
			}

			a.Children = widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("This page has not been implemented yet...", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})).Children
			widget.Refresh(a)
		}

		var currentPage = 0

		for {
			// Load contents to page when user is on the page.
			// Update only when the user is on the page.
			if menu.tabs.CurrentTabIndex() == 0 {
				if currentPage != 0 {
					history = historyPageData{}
					overviewPage(wallet, fyneApp)
					resetPages(0, window)
					currentPage = 0
				}
				// Todo: Remove overviewPageUpdate when TxNofier is implemented.
				overviewPageUpdates(wallet)
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
					overview = overviewPageData{}
					history = historyPageData{}
					notImplemented(2)
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

			time.Sleep(time.Millisecond * 50)
		}
	}()

	// Where peerConn and blkHeight are the realtime status texts.
	status := widget.NewVBox(menu.peerConn, menu.blkHeight)
	data := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, status, menu.tabs, nil), menu.tabs, status)

	return data
}

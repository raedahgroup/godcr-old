package pages

import (
	"context"
	"time"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/app/wallet"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type MenuPageStruct struct {
	tabMenu *widget.TabContainer
}

func pageNotImplemented() fyne.CanvasObject {
	label := widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	return widget.NewHBox(widgets.NewHSpacer(10), label)
}

func (menu *MenuPageStruct) MenuPage(ctx context.Context, wallet wallet.Wallet, window fyne.Window) {
	menu.tabMenu = widget.NewTabContainer(
		widget.NewTabItemWithIcon("Overview", overviewIcon, pageNotImplemented()),
		widget.NewTabItemWithIcon("History", historyIcon, widget.NewHBox()),
		widget.NewTabItemWithIcon("Send", sendIcon, widget.NewHBox()),
		widget.NewTabItemWithIcon("Receive", receiveIcon, widget.NewHBox()),
		widget.NewTabItemWithIcon("Accounts", accountIcon, widget.NewHBox()),
		widget.NewTabItemWithIcon("Staking", stakingIcon, widget.NewHBox()),
		widget.NewTabItemWithIcon("More", moreIcon, widget.NewHBox()),
		widget.NewTabItemWithIcon("Exit", exitIcon, widget.NewHBox()))
	menu.tabMenu.SetTabLocation(widget.TabLocationLeading)

	window.SetContent(menu.tabMenu)

	go func() {
		var currentTabIndex = -1

		for {
			if menu.tabMenu.CurrentTabIndex() == currentTabIndex {
				time.Sleep(50 * time.Millisecond)
				continue
			}

			// clear previous tab content
			previousTabIndex := currentTabIndex
			if previousTabIndex >= 0 {
				if previousPageBox, ok := menu.tabMenu.Items[previousTabIndex].Content.(*widget.Box); ok {
					previousPageBox.Children = widget.NewHBox().Children
					widget.Refresh(previousPageBox)
				}
			}

			currentTabIndex = menu.tabMenu.CurrentTabIndex()
			var newPageContent fyne.CanvasObject

			switch currentTabIndex {
			case 0:
				newPageContent = pageNotImplemented()
			case 1:
				newPageContent = pageNotImplemented()
			case 2:
				newPageContent = pageNotImplemented()
			case 3:
				newPageContent = pageNotImplemented()
			case 4:
				newPageContent = pageNotImplemented()
			case 5:
				newPageContent = pageNotImplemented()
			}

			if activePageBox, ok := menu.tabMenu.Items[currentTabIndex].Content.(*widget.Box); ok {
				activePageBox.Children = []fyne.CanvasObject{newPageContent}
				widget.Refresh(activePageBox)
			}
		}
	}()
}

package fyne

import (
	"fmt"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/pages"
)

func (app *fyneApp) displayLaunchErrorAndExit(errorMessage string) {
	app.window.SetContent(widget.NewVBox(
		widget.NewLabelWithStyle(errorMessage, fyne.TextAlignCenter, fyne.TextStyle{}),

		widget.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Exit", app.window.Close), // closing the window will trigger app.tearDown()
			layout.NewSpacer(),
		),
	))

	app.window.ShowAndRun()
	app.tearDown()
}

func (app *fyneApp) displayMainWindow() {
	app.setupNavigationMenu()
	app.tabMenu.SelectTabIndex(0)
	app.window.SetContent(app.tabMenu)

	app.window.ShowAndRun()
	app.tearDown()
}

func (app *fyneApp) setupNavigationMenu() {
	icons, err := getIcons(overviewIcon, historyIcon, sendIcon, receiveIcon, accountsIcon, stakeIcon)
	if err != nil {
		app.displayLaunchErrorAndExit(fmt.Sprintf("An error occured while loading app icons: %s", err))
	}
	app.tabMenu = widget.NewTabContainer(
		widget.NewTabItemWithIcon("Overview", icons[overviewIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("History", icons[historyIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Send", icons[sendIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Receive", icons[receiveIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Accounts", icons[accountsIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Staking", icons[stakeIcon], widget.NewHBox()),
	)
	app.tabMenu.SetTabLocation(widget.TabLocationLeading)

	go func() {
		var currentTabIndex = -1

		for {
			if app.tabMenu.CurrentTabIndex() == currentTabIndex {
				time.Sleep(50 * time.Millisecond)
				continue
			}

			// clear previous tab content
			previousTabIndex := currentTabIndex
			if previousTabIndex >= 0 {
				if previousPageBox, ok := app.tabMenu.Items[previousTabIndex].Content.(*widget.Box); ok {
					previousPageBox.Children = widget.NewHBox().Children
					widget.Refresh(previousPageBox)
				}
			}

			currentTabIndex = app.tabMenu.CurrentTabIndex()
			var newPageContent fyne.CanvasObject

			switch currentTabIndex {
			case 0:
				newPageContent = pages.OverviewPageContent()
			case 1:
				newPageContent = pages.HistoryPageContent()
			case 2:
				newPageContent = pages.SendPageContent()
			case 3:
				newPageContent = pages.ReceivePageContent()
			case 4:
				newPageContent = pages.AccountsPageContent()
			case 5:
				newPageContent = pages.StakingPageContent()
			}

			if activePageBox, ok := app.tabMenu.Items[currentTabIndex].Content.(*widget.Box); ok {
				activePageBox.Children = []fyne.CanvasObject{newPageContent}
				widget.Refresh(activePageBox)
			}
		}
	}()
}

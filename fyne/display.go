package fyne

import (
	"fmt"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/pages"
	"github.com/raedahgroup/godcr/fyne/utils"
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
	icons, err := utils.GetIcons(utils.OverviewIcon, utils.HistoryIcon, utils.SendIcon, utils.ReceiveIcon, utils.AccountsIcon, utils.StakeIcon)
	if err != nil {
		app.displayLaunchErrorAndExit(fmt.Sprintf("An error occured while loading app icons: %s", err))
	}
	app.tabMenu = widget.NewTabContainer(
		widget.NewTabItemWithIcon("Overview", icons[utils.OverviewIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("History", icons[utils.HistoryIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Send", icons[utils.SendIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Receive", icons[utils.ReceiveIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Accounts", icons[utils.AccountsIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Staking", icons[utils.StakeIcon], widget.NewHBox()),
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
				newPageContent = pages.ReceivePageContent(app.dcrlw, app.window, app.tabMenu)
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

package pages

import (
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/decred/slog"
	"github.com/raedahgroup/dcrlibwallet"
)

type AppInterface struct {
	Log            slog.Logger
	Dcrlw          *dcrlibwallet.LibWallet
	Window         fyne.Window
	AppDisplayName string

	tabMenu *widget.TabContainer
}

// DisplayLaunchErrorAndExit displays the error message to users.
func (app *AppInterface) DisplayLaunchErrorAndExit(errorMessage string) fyne.CanvasObject {
	return widget.NewVBox(
		widget.NewLabelWithStyle(errorMessage, fyne.TextAlignCenter, fyne.TextStyle{}),

		widget.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Exit", app.Window.Close), // closing the window will trigger app.tearDown()
			layout.NewSpacer(),
		))
}

func (app *AppInterface) MenuPage() {
	icons, err := getIcons(overviewIcon, historyIcon, sendIcon, receiveIcon, accountsIcon, stakeIcon)
	if err != nil {
		app.Window.SetContent(app.DisplayLaunchErrorAndExit(err.Error()))
		return
	}

	app.tabMenu = widget.NewTabContainer(
		widget.NewTabItemWithIcon("Overview", icons[overviewIcon], overviewPageContent()),
		widget.NewTabItemWithIcon("History", icons[historyIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Send", icons[sendIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Receive", icons[receiveIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Accounts", icons[accountsIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Staking", icons[stakeIcon], widget.NewHBox()),
	)
	app.tabMenu.SetTabLocation(widget.TabLocationLeading)

	app.Window.SetContent(app.tabMenu)

	go func() {
		var currentTabIndex = 0

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
				newPageContent = overviewPageContent()
			case 1:
				newPageContent = historyPageContent()
			case 2:
				newPageContent = sendPageContent()
			case 3:
				newPageContent = receivePageContent(app.Dcrlw, app.Window)
			case 4:
				newPageContent = accountsPageContent()
			case 5:
				newPageContent = stakingPageContent()
			}

			if activePageBox, ok := app.tabMenu.Items[currentTabIndex].Content.(*widget.Box); ok {
				activePageBox.Children = []fyne.CanvasObject{newPageContent}
				widget.Refresh(activePageBox)
			}
		}
	}()
}

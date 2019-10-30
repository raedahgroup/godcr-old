package pages

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/decred/slog"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/assets"
)

type AppInterface struct {
	Log            slog.Logger
	Dcrlw          *dcrlibwallet.LibWallet
	Window         fyne.Window
	AppDisplayName string

	tabMenu *widget.TabContainer
}

// DisplayLaunchErrorAndExit should only be used if ShowAndRun has already been called and it is not been return to a tabitem.
func (app *AppInterface) DisplayLaunchErrorAndExit(errorMessage string) {
	app.Window.SetContent(widget.NewVBox(
		widget.NewLabelWithStyle(errorMessage, fyne.TextAlignCenter, fyne.TextStyle{}),

		widget.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Exit", app.tearDown), // closing the window will trigger app.tearDown()
			layout.NewSpacer(),
		),
	))
	app.Window.ShowAndRun()
	app.tearDown()
	os.Exit(1)
}

func (app *AppInterface) displayErrorPage(errorMessage string) fyne.CanvasObject {
	return widget.NewVBox(
		widget.NewLabelWithStyle(errorMessage, fyne.TextAlignCenter, fyne.TextStyle{}),

		widget.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Exit", func() {
				app.tearDown()
				os.Exit(1)
			}),
			layout.NewSpacer(),
		),
	)
}

func (app *AppInterface) DisplayMainWindow() {
	app.setupNavigationMenu()
	app.Window.SetContent(app.tabMenu)
	app.Window.CenterOnScreen()
	app.Window.ShowAndRun()
	app.tearDown()
}

func (app *AppInterface) setupNavigationMenu() {
	icons, err := assets.Get(assets.OverviewIcon, assets.HistoryIcon, assets.SendIcon,
		assets.ReceiveIcon, assets.AccountsIcon, assets.StakingIcon)
	if err != nil {
		app.DisplayLaunchErrorAndExit(fmt.Sprintf("An error occured while loading app icons: %s", err))
		return
	}

	app.tabMenu = widget.NewTabContainer(
		widget.NewTabItemWithIcon("Overview", icons[assets.OverviewIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("History", icons[assets.HistoryIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Send", icons[assets.SendIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Receive", icons[assets.ReceiveIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Accounts", icons[assets.AccountsIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Staking", icons[assets.StakingIcon], widget.NewHBox()),
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
				newPageContent = overviewPageContent()
			case 1:
				newPageContent = historyPageContent()
			case 2:
				newPageContent = sendPageContent()
			case 3:
				newPageContent = receivePageContent()
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

	err = app.Dcrlw.SpvSync("") // todo dcrlibwallet should ideally read this parameter from config
	if err != nil {
		errorMessage := fmt.Sprintf("Spv sync attempt failed: %v", err)
		app.Log.Errorf(errorMessage)
		app.DisplayLaunchErrorAndExit(errorMessage)
		return
	}
}

func (app *AppInterface) tearDown() {
	if app.Dcrlw != nil {
		app.Dcrlw.Shutdown()
	}
	fyne.CurrentApp().Quit()
}

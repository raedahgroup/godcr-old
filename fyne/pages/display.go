package pages

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"


	"github.com/decred/slog"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
)

type AppInterface struct {
	Log            slog.Logger
	MultiWallet    *dcrlibwallet.MultiWallet
	Window         fyne.Window
	AppDisplayName string

	tabMenu *widget.TabContainer
}

func (app *AppInterface) DisplayLaunchErrorAndExit(errorMessage string) {
	app.Window.SetContent(widget.NewVBox(
		widget.NewLabelWithStyle(errorMessage, fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Exit", app.Window.Close), // closing the window will trigger app.tearDown()
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
			widget.NewButton("Exit", app.Window.Close), // closing the window will trigger app.tearDown()
			layout.NewSpacer(),
		),
	)
}

func (app *AppInterface) DisplayMainWindow() {
	app.setupNavigationMenu()
	app.Window.SetContent(app.tabMenu)
	app.Window.SetFixedSize(true)
	app.Window.CenterOnScreen()
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	go overviewHandler.PreserveSyncSteps()
	app.Window.ShowAndRun()
	app.tearDown()
}

func (app *AppInterface) setupNavigationMenu() {
	icons, err := assets.GetIcons(assets.OverviewIcon, assets.HistoryIcon, assets.SendIcon,
		assets.ReceiveIcon, assets.AccountsIcon, assets.StakeIcon)

	if err != nil {
		app.DisplayLaunchErrorAndExit(fmt.Sprintf("An error occured while loading app icons: %s", err))
		return
	}

	app.tabMenu = widget.NewTabContainer(
		widget.NewTabItemWithIcon("Overview", icons[assets.OverviewIcon], overviewPageContent(app)),
		widget.NewTabItemWithIcon("History", icons[assets.HistoryIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Send", icons[assets.SendIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Receive", icons[assets.ReceiveIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Accounts", icons[assets.AccountsIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Staking", icons[assets.StakeIcon], widget.NewHBox()),
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
				newPageContent = overviewPageContent(app)
			case 1:
				newPageContent = historyPageContent()
			case 2:
				newPageContent = sendPageContent(app.MultiWallet, app.Window)
			case 3:
				newPageContent = receivePageContent(app.MultiWallet, app.Window)
			case 4:
				newPageContent = accountsPageContent()
			case 5:
				newPageContent = stakingPageContent()
			}

			if activePageBox, ok := app.tabMenu.Items[currentTabIndex].Content.(*widget.Box); ok {
				activePageBox.Children = []fyne.CanvasObject{newPageContent}
				widget.Refresh(activePageBox)
				app.Window.Resize(app.tabMenu.MinSize().Union(newPageContent.MinSize()))
			}
		}
	}()

	err = app.MultiWallet.SpvSync()
	if err != nil {
		errorMessage := fmt.Sprintf("Spv sync attempt failed: %v", err)
		app.Log.Errorf(errorMessage)
		app.DisplayLaunchErrorAndExit(errorMessage)
		return
	}

	app.walletNotificationListener()
}

func (app *AppInterface) tearDown() {
	if app.MultiWallet != nil {
		app.MultiWallet.Shutdown()
	}
	fyne.CurrentApp().Quit()
}

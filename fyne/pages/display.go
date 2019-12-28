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
)

type AppInterface struct {
	Log            slog.Logger
	MultiWallet    *dcrlibwallet.MultiWallet
	Window         fyne.Window
	AppDisplayName string

	tabMenu  *widget.TabContainer
	handlers *pageHandlers
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
	// go overview.overviewHandler.PreserveSyncSteps()
	app.Window.ShowAndRun()
	app.tearDown()
}

func (app *AppInterface) setupNavigationMenu() {
	var err error
	p := initiatePages(app.MultiWallet, app.tabMenu, app.Window)

	app.handlers = p.pageHandlers()
	app.tabMenu, err = p.NewTab()
	if err != nil {
		app.DisplayLaunchErrorAndExit(fmt.Sprintf("An error occured while loading app icons: %s", err))
		return
	}

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
				newPageContent = p.overviewPage()
			case 1:
				newPageContent = p.historyPage()
			case 2:
				newPageContent = p.sendPage()
			case 3:
				newPageContent = p.receivePage()
			case 4:
				newPageContent = p.accountsPage()
			case 5:
				newPageContent = p.stakingPage()
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

package fyne

import (
	"fmt"

	"github.com/decred/slog"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/pages"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

type fyneApp struct {
	log     slog.Logger
	dcrlw   *dcrlibwallet.LibWallet
	window  fyne.Window
	tabMenu *widget.TabContainer
}

func LaunchUserInterface(appDisplayName, defaultAppDataDir, netType string) {
	fyne.SetCurrentApp(app.New())

	appInstance := &fyneApp{
		window: fyne.CurrentApp().NewWindow(appDisplayName),
	}

	appInstance.startUp(defaultAppDataDir, netType)
}

func (app *fyneApp) startUp(defaultAppDataDir, netType string) {
	var err error
	app.log, err = dcrlibwallet.RegisterLogger("FYNE")
	if err != nil {
		app.displayLaunchErrorAndExit(fmt.Sprintf("Logger setup error: %v", err))
		return
	}

	app.dcrlw, err = dcrlibwallet.NewLibWallet(defaultAppDataDir, "", netType)
	if err != nil {
		errorMessage := fmt.Sprintf("Initialization error: %v", err)
		app.log.Errorf(errorMessage)
		app.displayLaunchErrorAndExit(errorMessage)
		return
	}

	walletExists, err := app.dcrlw.WalletExists()
	if err != nil {
		errorMessage := fmt.Sprintf("Error checking if wallet db exists: %v", err)
		app.log.Errorf(errorMessage)
		app.displayLaunchErrorAndExit(errorMessage)
		return
	}

	if !walletExists {
		pages.ShowCreateAndRestoreWalletPage(app.dcrlw, app.window)
		return
	}

	// todo check settings.db to see if pub pass is configured and request from user
	// pass nil to use default pub pass
	err = app.dcrlw.OpenWallet(nil)
	if err != nil {
		errorMessage := fmt.Sprintf("Error opening wallet db: %v", err)
		app.log.Errorf(errorMessage)
		app.displayLaunchErrorAndExit(errorMessage)
		return
	}

	err = app.dcrlw.SpvSync("") // todo dcrlibwallet should ideally read this parameter from config
	if err != nil {
		errorMessage := fmt.Sprintf("Spv sync attempt failed: %v", err)
		app.log.Errorf(errorMessage)
		app.displayLaunchErrorAndExit(errorMessage)
		return
	}

	app.displayMainWindow()
}

func (app *fyneApp) tearDown() {
	if app.dcrlw != nil {
		app.dcrlw.Shutdown()
	}

	fyne.CurrentApp().Quit()
}

/*
go func() {
		var currentTabIndex = -1

		for {
			if app.tabMenu.CurrentTabIndex() == currentTabIndex {
				time.Sleep(50 * time.Millisecond)
				return
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
*/

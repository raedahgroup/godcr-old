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

func LaunchUserInterface(appDisplayName, appDataDir, netType string) {
	fyne.SetCurrentApp(app.New())

	appInstance := &fyneApp{
		window: fyne.CurrentApp().NewWindow(appDisplayName),
	}

	appInstance.startUp(appDataDir, netType)
}

func (app *fyneApp) startUp(appDataDir, netType string) {
	var err error
	app.log, err = dcrlibwallet.RegisterLogger("FYNE")
	if err != nil {
		app.displayLaunchErrorAndExit(fmt.Sprintf("Logger setup error: %v", err))
		return
	}

	app.dcrlw, err = dcrlibwallet.NewLibWallet(appDataDir, "", netType)
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
		app.setupNavigationMenu()
		pages.ShowCreateAndRestoreWalletPage(app.dcrlw, app.window, app.tabMenu, app.log)
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

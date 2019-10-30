package fyne

import (
	"fmt"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/pages"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
)

func LaunchUserInterface(appDisplayName, appDataDir, netType string) {
	fyne.SetCurrentApp(app.New())

	var app = pages.AppInterface{
		Window:         fyne.CurrentApp().NewWindow(appDisplayName),
		AppDisplayName: appDisplayName,
	}

	var err error
	app.Log, err = dcrlibwallet.RegisterLogger("FYNE")
	if err != nil {
		app.DisplayLaunchErrorAndExit(fmt.Sprintf("Logger setup error: %v", err))
		return
	}

	app.Dcrlw, err = dcrlibwallet.NewLibWallet(appDataDir, "", netType)
	if err != nil {
		errorMessage := fmt.Sprintf("Initialization error: %v", err)
		app.Log.Errorf(errorMessage)
		app.DisplayLaunchErrorAndExit(errorMessage)
		return
	}

	walletExists, err := app.Dcrlw.WalletExists()
	if err != nil {
		errorMessage := fmt.Sprintf("Error checking if wallet db exists: %v", err)
		app.Log.Errorf(errorMessage)
		app.DisplayLaunchErrorAndExit(errorMessage)
		return
	}

	if !walletExists {
		app.ShowCreateAndRestoreWalletPage()
		return
	}

	// todo check settings.db to see if pub pass is configured and request from user
	// pass nil to use default pub pass
	err = app.Dcrlw.OpenWallet(nil)
	if err != nil {
		errorMessage := fmt.Sprintf("Error opening wallet db: %v", err)
		app.Log.Errorf(errorMessage)
		app.DisplayLaunchErrorAndExit(errorMessage)
		return
	}

	app.DisplayMainWindow()
}

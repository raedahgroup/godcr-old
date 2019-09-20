package fyne

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/decred/slog"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/pages"
)

type fyneApp struct {
	log    slog.Logger
	dcrlw  *dcrlibwallet.LibWallet
	window fyne.Window
}

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

func LaunchUserInterface(appDisplayName, appDataDir, netType string) {
	fyne.SetCurrentApp(app.New())

	appInstance := &fyneApp{
		window: fyne.CurrentApp().NewWindow(appDisplayName),
	}

	appInstance.startUp(appDataDir, appDisplayName, netType)
}

func (app *fyneApp) startUp(appDataDir, appDisplayName, netType string) {
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

	var menu = pages.AppInterface{
		Dcrlw:          app.dcrlw,
		Window:         app.window,
		Log:            app.log,
		AppDisplayName: appDisplayName,
	}

	if !walletExists {
		menu.ShowCreateAndRestoreWalletPage()
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

	err = app.dcrlw.SpvSync("")
	if err != nil {
		errorMessage := fmt.Sprintf("Spv sync attempt failed: %v", err)
		app.log.Errorf(errorMessage)
		app.displayLaunchErrorAndExit(errorMessage)
		return
	}

	menu.MenuPage()
	menu.Window.CenterOnScreen()
	menu.Window.Resize(fyne.NewSize(500, 500))
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	menu.Window.ShowAndRun()
	app.tearDown()
}

func (app *fyneApp) tearDown() {
	if app.dcrlw != nil {
		app.dcrlw.Shutdown()
	}

	fyne.CurrentApp().Quit()
}

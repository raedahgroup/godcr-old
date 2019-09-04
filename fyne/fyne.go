package fyne

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"

	"github.com/raedahgroup/dcrlibwallet"
)

type fyneApp struct {
	defaultAppDataDir string
	netType string
	dcrlw *dcrlibwallet.LibWallet
	window fyne.Window
}

func InitializeUserInterface(appDisplayName, defaultAppDataDir, netType string) *fyneApp {
	// set app instance to be accessed subsequently as fyne.CurrentApp()
	fyne.SetCurrentApp(app.New())

	return &fyneApp{
		defaultAppDataDir: defaultAppDataDir,
		netType: netType,
		window: fyne.CurrentApp().NewWindow(appDisplayName),
	}
}

func (app *fyneApp) Launch() {
	dcrlw, err := dcrlibwallet.NewLibWallet(app.defaultAppDataDir, "", app.netType)
	if err != nil {
		// todo: display prelaunch error
		return
	}

	app.dcrlw = dcrlw
}

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

func LaunchUserInterface(appDisplayName, defaultAppDataDir, netType string) {
	// set app instance to be accessed subsequently as fyne.CurrentApp()
	fyne.SetCurrentApp(app.New())

	f := &fyneApp{
		defaultAppDataDir: defaultAppDataDir,
		netType: netType,
		window: fyne.CurrentApp().NewWindow(appDisplayName),
	}

	f.Launch()
}

func (app *fyneApp) Launch() {
	dcrlw, err := dcrlibwallet.NewLibWallet(app.defaultAppDataDir, "", app.netType)
	if err != nil {
		// todo: display prelaunch error
		return
	}

	app.dcrlw = dcrlw
}

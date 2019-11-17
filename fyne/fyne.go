package fyne

import (
	"fmt"
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/app"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/pages"
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

	app.MultiWallet, err = dcrlibwallet.NewMultiWallet(appDataDir, "", netType)
	if err != nil {
		errorMessage := fmt.Sprintf("Initialization error: %v", err)
		app.Log.Errorf(errorMessage)
		app.DisplayLaunchErrorAndExit(errorMessage)
		return
	}

	walletCount := app.MultiWallet.LoadedWalletsCount()

	if walletCount == 0 {
		app.ShowCreateAndRestoreWalletPage()
		return
	}

	// todo check settings.db to see if pub pass is configured and request from user
	// pass nil to use default pub pass
	err = app.MultiWallet.OpenWallets(nil)
	if err != nil {
		errorMessage := fmt.Sprintf("Error opening wallet db: %v", err)
		app.Log.Errorf(errorMessage)
		app.DisplayLaunchErrorAndExit(errorMessage)
		return
	}

	app.Wallets = make([]*dcrlibwallet.Wallet, walletCount)
	openedWallets := app.MultiWallet.OpenedWalletIDsRaw()
	sort.Ints(openedWallets)
	for walletIndex, walletID := range openedWallets {
		app.Wallets[walletIndex] = app.MultiWallet.WalletWithID(walletID)
	}

	app.DisplayMainWindow()
}

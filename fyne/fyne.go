package fyne

import (
	"context"
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/app"

	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/wallet"
	"github.com/raedahgroup/godcr/fyne/pages"
)

type fyneApp struct {
	ctx    context.Context
	cfg    *config.Config
	wallet wallet.Wallet
}

func InitializeUserInterface() *fyneApp {
	// set app instance to be accessed subsequently as fyne.CurrentApp()
	fyne.SetCurrentApp(app.New())
	return &fyneApp{}
}

func (app *fyneApp) DisplayPreLaunchError(errorMessage string) {
	// todo: show an error message dialog
}

func (app *fyneApp) LaunchApp(ctx context.Context, cfg *config.Config, wallet wallet.Wallet) {
	app.ctx = ctx
	app.cfg = cfg
	app.wallet = wallet

	walletExists, err := app.wallet.WalletExists()
	if err != nil {
		errorMessage := fmt.Sprintf("Error checking if wallet db exists: %v", err)
		log.Errorf(errorMessage)
		app.DisplayPreLaunchError(errorMessage)
		return
	}

	if !walletExists {
		pages.ShowWelcomePage(app.wallet)
		return
	}

	err = app.wallet.OpenWallet(app.ctx, "")
	if err != nil {
		errorMessage := fmt.Sprintf("Error opening wallet db: %v", err)
		log.Errorf(errorMessage)
		app.DisplayPreLaunchError(errorMessage)
		return
	}

	err = app.wallet.SpvSync(false, nil)
	if err != nil {
		errorMessage := fmt.Sprintf("Spv sync attempt failed: %v", err)
		log.Errorf(errorMessage)
		app.DisplayPreLaunchError(errorMessage)
		return
	}

	// todo: display overview page (include sync progress UI elements)
	// todo: register sync progress listener on overview page to update sync progress views
}

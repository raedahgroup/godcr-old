package fyne

import (
	"context"

	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/wallet"
)

type fyneApp struct {
	ctx    context.Context
	cfg    *config.Config
	wallet wallet.Wallet
}

func InitializeUserInterface() *fyneApp {
	return &fyneApp{}
}

func (app *fyneApp) HandlePreLaunchError(err error) {
	// todo: show an error message dialog
}

func (app *fyneApp) LaunchApp(ctx context.Context, cfg *config.Config, wallet wallet.Wallet) {
	app.ctx = ctx
	app.cfg = cfg
	app.wallet = wallet
}

package fyne

import (
	"context"
	"fmt"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/wallet"
	"github.com/raedahgroup/godcr/fyne/pages"
)

type fyneApp struct {
	ctx    context.Context
	cfg    *config.Config
	wallet wallet.Wallet
	window fyne.Window
}

func InitializeUserInterface() *fyneApp {
	// set app instance to be accessed subsequently as fyne.CurrentApp()
	fyne.SetCurrentApp(app.New())
	window := fyne.CurrentApp().NewWindow(godcrApp.DisplayName)
	return &fyneApp{window: window}
}

func (app *fyneApp) DisplayPreLaunchError(errorMessage string) {
	app.window.SetContent(widget.NewVBox(
		widget.NewLabelWithStyle(errorMessage, fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewHBox(layout.NewSpacer(), widget.NewButton("Exit", func() {
			app.window.Close()
			fyne.CurrentApp().Quit()
			os.Exit(1)
		}), layout.NewSpacer())))

	app.window.ShowAndRun()
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
		pages.ShowCreateAndRestoreWalletPage(app.wallet, app.window, ctx)
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

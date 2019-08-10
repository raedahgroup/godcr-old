package fyne

import (
	"context"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/fyne/pages"
)

//todo: implement better error logger, probably log to a file
func LaunchFyne(ctx context.Context, walletMiddleware godcrApp.WalletMiddleware, config config.Settings) {
	fyneApp := app.New()
	window := fyneApp.NewWindow(godcrApp.DisplayName)
	sync := pages.ShowSyncWindow(ctx, walletMiddleware, config, window, fyneApp)
	window.SetContent(sync)
	window.Resize(fyne.NewSize(1000, 500))
	window.CenterOnScreen()
	window.ShowAndRun()
}

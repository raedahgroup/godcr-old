package fyne

import (
	"context"

	"fyne.io/fyne"

	"fyne.io/fyne/app"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/pages"
)

func LaunchFyne(ctx context.Context, walletMiddleware godcrApp.WalletMiddleware) {
	fyneApp := app.New()
	window := fyneApp.NewWindow(godcrApp.DisplayName)
	sync := pages.ShowSyncWindow(walletMiddleware, window, fyneApp)
	window.Resize(fyne.NewSize(1000, 500))
	window.SetContent(sync)
	window.CenterOnScreen()
	window.ShowAndRun()
}

package fyne

import (
	"context"

	"fyne.io/fyne/app"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/pages"
)

func LaunchFyne(ctx context.Context, walletMiddleware godcrApp.WalletMiddleware) {
	fyneApp := app.New()
	window := fyneApp.NewWindow(godcrApp.DisplayName)
	window.SetContent(pages.Menu(pages.ShowSyncWindow(walletMiddleware, window, fyneApp), window, fyneApp))
	window.ShowAndRun()

}

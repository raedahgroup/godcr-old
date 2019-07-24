package fyne

import (
	"context"

	"fyne.io/fyne/app"
	godcrApp "github.com/raedahgroup/godcr/app"
)

func LaunchFyne(ctx context.Context, walletMiddleware godcrApp.WalletMiddleware) {
	fyneApp := app.New()
	window := fyneApp.NewWindow(godcrApp.DisplayName)
	window.SetContent(ShowSyncWindow(walletMiddleware, window, fyneApp))
	window.ShowAndRun()

}

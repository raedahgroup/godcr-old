package terminal

import (
	"context"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/pages"
	"github.com/rivo/tview"
)

// todo the ctx variable should be stored somewhere for as long as this terminal app is open
// it will be necessary for use in some wallet operations
func StartTerminalApp(_ context.Context, walletMiddleware app.WalletMiddleware) error {
	tviewApp := tview.NewApplication()

	walletExists, err := walletMiddleware.WalletExists()
	if err != nil {
		return err
	}
	if walletExists {
		pages.LaunchSyncPage(tviewApp, walletMiddleware)
	} else {
		tviewApp.SetRoot(pages.CreateWalletPage(tviewApp, walletMiddleware), true)
	}

	// `Run` blocks until app.Stop() is called before returning
	return tviewApp.Run()
}

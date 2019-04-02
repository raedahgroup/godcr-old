package terminal

import (
	"context"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/pages"
	"github.com/rivo/tview"
)

func StartTerminalApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	tviewApp := tview.NewApplication()

	walletExists, err := walletMiddleware.OpenWalletIfExist(ctx)
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

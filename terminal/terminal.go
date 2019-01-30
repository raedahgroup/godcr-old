package terminal

import (
	"context"
	"fmt"
	"strings"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/cli/walletloader"
	"github.com/rivo/tview"
)

func StartTerminalApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	tviewApp := tview.NewApplication()
	list := tview.NewList().
		AddItem("Balance", "", 'b', nil).
		AddItem("Receive", "", 'r', nil).
		AddItem("Send", "", 's', nil).
		AddItem("History", "", 'h', nil).
		AddItem("Exit", "", 'q', func() {
			tviewApp.Stop()
		})
	list.SetBorder(true).SetTitle(fmt.Sprintf("%s Terminal", strings.ToUpper(app.Name)))

	err := syncBlockChain(ctx, walletMiddleware)
	if err != nil {
		fmt.Println(err)
	}

	// `Run` blocks until app.Stop() is called before returning
	return tviewApp.SetRoot(list, true).SetFocus(list).Run()
}

func syncBlockChain(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	_, err := walletloader.OpenWallet(ctx, walletMiddleware)
	if err != nil {
		return err
	}

	return walletloader.SyncBlockChain(ctx, walletMiddleware)
}

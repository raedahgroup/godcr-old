package terminal

import (
	"context"
	"fmt"
	"strings"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/cli/walletloader"
	"github.com/rivo/tview"
)

var tviewApp *tview.Application

func StartTerminalApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	tviewApp = tview.NewApplication()
	
	// Grid Layout Structure for screens
	tviewGrid := GridLayout()

	err := syncBlockChain(ctx, walletMiddleware)
	if err != nil{
		fmt.Println(err)
	}

	// `Run` blocks until app.Stop() is called before returning
	return tviewApp.SetRoot(tviewGrid, true).SetFocus(tviewGrid).Run()
}

func GridLayout() tview.Primitive {

	PagesColumn := func(text string) tview.Primitive {
		return tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(text)
	}

	//List of menu that will display at the left hand
	MenuColumn := func() tview.Primitive {
		return tview.NewList().
			AddItem("Balance", "", 'b', nil).
			AddItem("Receive", "", 'r', nil).
			AddItem("Send", "", 's', nil).
			AddItem("History", "", 'h', nil).
			AddItem("Exit", "", 'q', func() {
			tviewApp.Stop()
		})
	}

	newmenu := MenuColumn()
	newpages := PagesColumn("Pages")
	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("%s Terminal", strings.ToUpper(app.Name)))

	gridApp := tview.NewGrid().SetBorders(true).SetRows(3, 0).SetColumns(30, 0)

	// First Row With single Column
	gridApp.AddItem(header, 0, 0, 1, 3, 0, 0, false)

	// Second Row with two column
	gridApp.AddItem(newmenu, 1, 0, 1, 1, 0, 0, true).
			AddItem(newpages, 1, 1, 1, 2, 0, 0, false)

	return gridApp
}

func syncBlockChain(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	_, err := walletloader.OpenWallet(ctx, walletMiddleware)
	if err != nil {
		return err
	}

	return walletloader.SyncBlockChain(ctx, walletMiddleware)
}
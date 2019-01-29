package terminal

import (
	"context"
	"fmt"
	"strings"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/cli/walletloader"
	"github.com/rivo/tview"
	"github.com/raedahgroup/godcr/terminal/pages"
)

func StartTerminalApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	var tviewApp *tview.Application

	tviewApp = tview.NewApplication()

	// Terminal Layout Structure for screens
	terminalLayoutApp := terminalLayout(tviewApp)
	
	err := syncBlockChain(ctx, walletMiddleware)
	if err != nil {
		fmt.Println(err)
	}

	// `Run` blocks until app.Stop() is called before returning
	return tviewApp.SetRoot(terminalLayoutApp, true).SetFocus(terminalLayoutApp).Run()
}

func terminalLayout(tviewApp *tview.Application) tview.Primitive {
	var menuColumn *tview.List

	//The Default page onload
	currentPage := pages.BalancePage()

	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("%s Terminal", strings.ToUpper(app.Name)))
	
	//Creating the View for the Layout
	gridLayoutApp := tview.NewGrid().SetBorders(true).SetRows(3, 0).SetColumns(30, 0)

	//Controls the display for the right side column
	changePageColumn := func(t tview.Primitive) {
			currentPage = t
			gridLayoutApp.AddItem(currentPage, 1, 1, 1, 2, 0, 0, true)
	}

	menuColumn =  tview.NewList().
		AddItem("Balance", "", 'b', func() {
		    changePageColumn(pages.BalancePage())
		}).
		AddItem("Receive", "", 'r', func() {
		    changePageColumn(pages.ReceivePage())
		}).
		AddItem("Send", "", 's', func() {
		    changePageColumn(pages.SendPage(tviewApp, menuColumn))
		}).
		AddItem("History", "", 'h', func() {
		    changePageColumn(pages.HistoryPage())
		}).
		AddItem("Exit", "", 'q', func() {
		    tviewApp.Stop()
	})

	menuColumn.SetCurrentItem(0)


	// Layout for screens Header
	gridLayoutApp.AddItem(header, 0, 0, 1, 3, 0, 0, false)

	// Layout for screens with two column
	gridLayoutApp.AddItem(menuColumn, 1, 0, 1, 1, 0, 0, true)
	gridLayoutApp.AddItem(currentPage, 1, 1, 1, 2, 0, 0, true)

	return gridLayoutApp
}

func syncBlockChain(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	_, err := walletloader.OpenWallet(ctx, walletMiddleware)
	if err != nil {
		return err
	}

	return walletloader.SyncBlockChain(ctx, walletMiddleware)
}

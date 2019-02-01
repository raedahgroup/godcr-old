package terminal

import (
	"context"
	"fmt"
	"strings"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/cli/walletloader"
	"github.com/raedahgroup/godcr/terminal/pages"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
)

func StartTerminalApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	tviewApp := tview.NewApplication()
	// Terminal Layout Structure for screens
	layout := terminalLayout(tviewApp)

	err := syncBlockChain(ctx, walletMiddleware)
	if err != nil {
		fmt.Println(err)
	}
	// `Run` blocks until app.Stop() is called before returning
	return tviewApp.SetRoot(layout, true).SetFocus(layout).Run()
}

func terminalLayout(tviewApp *tview.Application) tview.Primitive {
	var menuColumn *tview.List

	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("%s Terminal", strings.ToUpper(app.Name)))
	header.SetBackgroundColor(tcell.NewRGBColor(41, 112, 255))
	//Creating the View for the Layout
	gridLayout := tview.NewGrid().SetBorders(false).SetRows(3, 0).SetColumns(30, 0)
	//Controls the display for the right side column
	changePageColumn := func(t tview.Primitive) {
		gridLayout.AddItem(t, 1, 1, 1, 1, 0, 0, true)
	}
	//Menu List of the Layout
	menuColumn = tview.NewList().
		AddItem("Balance", "", 'b', func() {
			changePageColumn(pages.BalancePage())
		}).
		AddItem("Receive", "", 'r', func() {
			changePageColumn(pages.ReceivePage())
		}).
		AddItem("Send", "", 's', func() {
			setFocus := tviewApp.SetFocus
			clearFocus := func() {
				tviewApp.SetFocus(menuColumn)
			}
			changePageColumn(pages.SendPage(setFocus, clearFocus))
		}).
		AddItem("History", "", 'h', func() {
			changePageColumn(pages.HistoryPage())
		}).
		AddItem("Exit", "", 'q', func() {
			tviewApp.Stop()
		})
	menuColumn.SetCurrentItem(0)
	menuColumn.SetBackgroundColor(tcell.NewRGBColor(0, 0, 51))
	// Layout for screens Header
	gridLayout.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	// Layout for screens with two column
	gridLayout.AddItem(menuColumn, 1, 0, 1, 1, 0, 0, true)
	gridLayout.AddItem(pages.BalancePage(), 1, 1, 1, 1, 0, 0, true)

	return gridLayout
}

func syncBlockChain(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	_, err := walletloader.OpenWallet(ctx, walletMiddleware)
	if err != nil {
		return err
	}

	return walletloader.SyncBlockChain(ctx, walletMiddleware)
}

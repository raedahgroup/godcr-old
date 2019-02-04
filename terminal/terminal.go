package terminal

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/cli/walletloader"
	"github.com/raedahgroup/godcr/terminal/pages"
	"github.com/rivo/tview"
)

func StartTerminalApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	tviewApp := tview.NewApplication()
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

	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("\n%s Terminal", strings.ToUpper(app.Name)))
	header.SetBackgroundColor(tcell.NewRGBColor(41, 112, 255))
	//Creating the View for the Layout
	gridLayout := tview.NewGrid().SetBorders(false).SetRows(3, 0).SetColumns(30, 0)
	//Controls the display for the right side column
	changePageColumn := func(t tview.Primitive) {
		gridLayout.AddItem(t, 1, 1, 1, 1, 0, 0, true)
		gridLayout.RemoveItem(t)
	}
	//Menu List of the Layout
	menuColumn = tview.NewList().
		AddItem("Overview", "", 'o', func() {
			changePageColumn(pages.BalancePage())
		}).
		AddItem("History", "", 'h', func() {
			changePageColumn(pages.HistoryPage())
		}).
		AddItem("Send", "", 's', func() {
			setFocus := tviewApp.SetFocus
			clearFocus := func() {
				tviewApp.SetFocus(menuColumn)
			}
			changePageColumn(pages.SendPage(setFocus, clearFocus))
		}).
		AddItem("Receive", "", 'r', func() {
			changePageColumn(pages.ReceivePage())
		}).
		AddItem("Stakeinfo", "", 'k', func() {
			changePageColumn(pages.ReceivePage())
		}).
		AddItem("Purchase Tickets", "", 't', func() {
			changePageColumn(pages.ReceivePage())
		}).
		AddItem("Account", "", 'a', nil).
		AddItem("Security", "", 'x', nil).
		AddItem("Quit", "", 'q', func() {
			tviewApp.Stop()
		})
	menuColumn.SetCurrentItem(0)
	menuColumn.SetShortcutColor(tcell.NewRGBColor(112, 203, 255))
	menuColumn.SetBorder(true)
	menuColumn.SetBorderColor(tcell.NewRGBColor(112, 203, 255))
	// Layout for screens Header
	gridLayout.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	// Layout for screens with two column
	gridLayout.AddItem(menuColumn, 1, 0, 1, 1, 0, 0, true)
	gridLayout.AddItem(pages.BalancePage(), 1, 1, 1, 1, 0, 0, true)
	gridLayout.SetBackgroundColor(tcell.ColorBlack)
	gridLayout.SetBordersColor(tcell.NewRGBColor(112, 203, 255))

	return gridLayout
}

func syncBlockChain(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	_, err := walletloader.OpenWallet(ctx, walletMiddleware)
	if err != nil {
		return err
	}

	return walletloader.SyncBlockChain(ctx, walletMiddleware)
}
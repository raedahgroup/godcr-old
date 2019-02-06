package terminal

import (
	"context"
	"fmt"
	"strings"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/pages"
	"github.com/rivo/tview"
)

func StartTerminalApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
<<<<<<< HEAD
	tviewApp := tview.NewApplication()
	layout := terminalLayout(tviewApp, walletMiddleware)

	// open wallet and start blockchain syncing in background
	walletExists, err := openWalletIfExist(ctx, walletMiddleware)
=======

	tviewApp := tview.NewApplication()
	layout := terminalLayout(tviewApp, walletMiddleware)
	
	err := syncBlockChain(ctx, walletMiddleware)
>>>>>>> implement the terminal UI history page
	if err != nil {
		return err
	}

	if walletExists {
		// `Run` blocks until app.Stop() is called before returning
		return tviewApp.SetRoot(layout, true).SetFocus(layout).Run()
	}

	var password, confPassword string
	form := tview.NewForm().
		AddPasswordField("Password", "", 10, '*', func(text string) {
			if len(text) == 0 {
				fmt.Println("password cannot be less than 4")
				return
			}
			password = text
		}).
		AddPasswordField("Password", "", 10, '*', func(text string) {
			confPassword = text
		}).
		AddButton("Create", func() {
			if password != confPassword {
				fmt.Println("password does not match")
				return
			}
			CreateWallet(ctx, password, walletMiddleware)
		}).
		AddButton("Quit", func() {
			tviewApp.Stop()
		})
	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignCenter).SetRect(30, 10, 40, 10)
	return tviewApp.SetRoot(form, false).Run()
}

func terminalLayout(tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {
	var menuColumn *tview.List

	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("%s Terminal", strings.ToUpper(app.Name)))
	//Creating the View for the Layout
	gridLayout := tview.NewGrid().SetBorders(true).SetRows(3, 0).SetColumns(30, 0)
	//Controls the display for the right side column
	changePageColumn := func(t tview.Primitive) {
		gridLayout.AddItem(t, 1, 1, 1, 1, 0, 0, true)
	}

	setFocus := tviewApp.SetFocus
	clearFocus := func() {
		tviewApp.SetFocus(menuColumn)
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
			changePageColumn(pages.SendPage(setFocus, clearFocus))
		}).
		AddItem("History", "", 'h', func() {
			changePageColumn(pages.HistoryPage(walletMiddleware))
		}).
		AddItem("Stakeinfo", "", 'k', func() {
			changePageColumn(pages.ReceivePage())
		}).
		AddItem("Purchase Tickets", "", 't', func() {
			changePageColumn(pages.ReceivePage())
		}).
		AddItem("Exit", "", 'q', func() {
			tviewApp.Stop()
		})
	menuColumn.SetCurrentItem(0)
	// Layout for screens Header
	gridLayout.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	// Layout for screens with two column
	gridLayout.AddItem(menuColumn, 1, 0, 1, 1, 0, 0, true)
	gridLayout.AddItem(pages.BalancePage(), 1, 1, 1, 1, 0, 0, true)

	return gridLayout
}

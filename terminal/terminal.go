package terminal

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/pages"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
)

func StartTerminalApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	tviewApp := tview.NewApplication()

	// open wallet and start blockchain syncing in background
	walletExists, err := openWalletIfExist(ctx, walletMiddleware)
	if err != nil {
		return err
	}
	if walletExists {
		// `Run` blocks until app.Stop() is called before returning
		layout := terminalLayout(tviewApp, walletMiddleware)
		return tviewApp.SetRoot(layout, true).SetFocus(layout).Run()
	}

	var password, confPassword string
	form := tview.NewForm().
		AddPasswordField("Password", "", 20, '*', func(text string) {
			if len(text) == 0 {
				fmt.Println("password cannot be less than 4")
				return
			}
			password = text
		}).
		AddPasswordField("Password", "", 20, '*', func(text string) {
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
	header.SetBackgroundColor(tcell.NewRGBColor(41, 112, 255))
	//Creating the View for the Layout
	gridLayout := tview.NewGrid().SetBorders(false).SetRows(3, 0).SetColumns(30, 0)
	//Controls the display for the right side column
	changePageColumn := func(t tview.Primitive) {
		gridLayout.AddItem(t, 1, 1, 1, 1, 0, 0, true)
		gridLayout.RemoveItem(t)
	}

	setFocus := tviewApp.SetFocus
	clearFocus := func() {
		tviewApp.SetFocus(menuColumn)
	}
	//Menu List of the Layout
	menuColumn = tview.NewList().
		AddItem("Balance", "", 'b', func() {
			changePageColumn(pages.BalancePage(walletMiddleware, setFocus, clearFocus))
		}).
		AddItem("Receive", "", 'r', func() {
			changePageColumn(pages.ReceivePage(walletMiddleware, setFocus, clearFocus))
		}).
		AddItem("Send", "", 's', func() {
			changePageColumn(pages.SendPage(setFocus, clearFocus))
		}).
		AddItem("History", "", 'h', func() {
			changePageColumn(pages.HistoryPage(walletMiddleware, setFocus, clearFocus))
		}).
		AddItem("Stakeinfo", "", 'k', func() {
			changePageColumn(pages.StakeinfoPage(walletMiddleware, setFocus, clearFocus))
		}).
		AddItem("Purchase Tickets", "", 't', func() {
			changePageColumn(pages.BalancePage())
		}).
		AddItem("Account", "", 'a', nil).
		AddItem("Security", "", 'x', nil).
		AddItem("Quit", "", 'q', func() {
			tviewApp.Stop()
		})
	menuColumn.SetCurrentItem(0)
	menuColumn.SetShortcutColor(tcell.NewRGBColor(112, 203, 255))
	menuColumn.SetBorder(true)
	menuColumn.SetBorderColor(pages.BorderColor)
	// Layout for screens Header
	gridLayout.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	// Layout for screens with two column
	gridLayout.AddItem(menuColumn, 1, 0, 1, 1, 0, 0, true)
	gridLayout.AddItem(pages.BalancePage(walletMiddleware, setFocus, clearFocus), 1, 1, 1, 1, 0, 0, true)
	gridLayout.SetBackgroundColor(tcell.ColorBlack)
	gridLayout.SetBordersColor(pages.BorderColor)

	return gridLayout
}

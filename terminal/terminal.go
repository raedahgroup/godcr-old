package terminal

import (
	"context"
	"fmt"
	"strings"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/pages"
	"github.com/rivo/tview"
)

func StartTerminalApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	tviewApp := tview.NewApplication()

	// open wallet and start blockchain syncing in background
	walletExists, err := openWalletIfExist(ctx, walletMiddleware)
	if err != nil {
		return err
	}
	if walletExists {
		err := SyncBlockChain(ctx, walletMiddleware)
		if err != nil {
			fmt.Println(err)
		}
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

	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("\n%s Terminal", strings.ToUpper(app.Name)))
	header.SetBackgroundColor(helpers.DecredColor)
	//Creating the View for the Layout
	gridLayout := tview.NewGrid().SetRows(3, 0).SetColumns(30, 0)
	//Controls the display for the right side column
	changePageColumn := func(t tview.Primitive) {
		gridLayout.AddItem(t, 1, 1, 1, 1, 0, 0, true)
	}

	setFocus := tviewApp.SetFocus

	menuColumn := tview.NewList()
	var page tview.Primitive
	clearFocus := func() {
		gridLayout.RemoveItem(page)
		tviewApp.SetFocus(menuColumn)
	}

	//Menu List of the Layout
	menuColumn.AddItem("Balance", "", 'b', func() {
		page = pages.BalancePage(walletMiddleware, setFocus, clearFocus)
		changePageColumn(page)
	}).SetSelectedFocusOnly(true)

	menuColumn.AddItem("Receive", "", 'r', func() {
		page = pages.ReceivePage(walletMiddleware, setFocus, clearFocus)
		changePageColumn(page)
	}).SetSelectedFocusOnly(true)

	menuColumn.AddItem("Send", "", 's', func() {
		page = pages.SendPage(setFocus, clearFocus)
		changePageColumn(page)
	}).SetSelectedFocusOnly(true)

	menuColumn.AddItem("History", "", 'h', func() {
		page = pages.HistoryPage(walletMiddleware, setFocus, clearFocus)
		changePageColumn(page)
	}).SetSelectedFocusOnly(true)

	menuColumn.AddItem("Stakeinfo", "", 'k', func() {
		page = pages.StakeinfoPage(walletMiddleware, setFocus, clearFocus)
		changePageColumn(page)
	}).SetSelectedFocusOnly(true)

	menuColumn.AddItem("Purchase Tickets", "", 't', func() {
		page = pages.PurchaseTicketsPage(walletMiddleware, setFocus, clearFocus)
		changePageColumn(page)
	}).SetSelectedFocusOnly(true)

	menuColumn.AddItem("Exit", "", 'q', func() {
		tviewApp.Stop()
	}).SetSelectedFocusOnly(true)

	menuColumn.SetCurrentItem(0)
	menuColumn.SetShortcutColor(helpers.DecredLightColor)
	menuColumn.SetBorder(true)
	menuColumn.SetBorderColor(helpers.DecredLightColor)
	// Layout for screens Header
	gridLayout.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	// Layout for screens with two column
	gridLayout.AddItem(menuColumn, 1, 0, 1, 1, 0, 0, true)
	gridLayout.SetBordersColor(helpers.DecredLightColor)

	return gridLayout
}

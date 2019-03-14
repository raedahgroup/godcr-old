package terminal

import (
	"context"
	"fmt"
	"strings"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/pages"
	"github.com/raedahgroup/godcr/terminal/primitives"
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
		return tviewApp.SetRoot(layout, true).Run()
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
	gridLayout := tview.NewGrid().
		SetRows(3, 0).
		SetColumns(30, 0).
		SetBordersColor(helpers.DecredLightColor)

	var activePage tview.Primitive

	// controls the display for the right side column
	changePageColumn := func(page tview.Primitive) {
		gridLayout.RemoveItem(activePage)
		activePage = page
		gridLayout.AddItem(activePage, 1, 1, 1, 1, 0, 0, true)
	}

	menuColumn := tview.NewList()
	clearFocus := func() {
		gridLayout.RemoveItem(activePage)
		tviewApp.SetFocus(menuColumn)
	}

	menuColumn.AddItem("Balance", "", 'b', func() {
		changePageColumn(pages.BalancePage(walletMiddleware, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Receive", "", 'r', func() {
		changePageColumn(pages.ReceivePage(walletMiddleware, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Send", "", 's', func() {
		changePageColumn(pages.SendPage(tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("History", "", 'h', func() {
		changePageColumn(pages.HistoryPage(walletMiddleware, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Staking", "", 'k', func() {
		changePageColumn(pages.StakingPage(walletMiddleware, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Exit", "", 'q', func() {
		tviewApp.Stop()
	})

	header := primitives.NewCenterAlignedTextView(fmt.Sprintf("\n%s Terminal", strings.ToUpper(app.Name)))
	header.SetBackgroundColor(helpers.DecredColor)
	gridLayout.AddItem(header, 0, 0, 1, 2, 0, 0, false)

	menuColumn.SetShortcutColor(helpers.DecredLightColor)
	menuColumn.SetBorder(true).SetBorderColor(helpers.DecredLightColor)
	gridLayout.AddItem(menuColumn, 1, 0, 1, 1, 0, 0, true)

	menuColumn.SetCurrentItem(0)
	changePageColumn(pages.BalancePage(walletMiddleware, tviewApp.SetFocus, clearFocus))

	return gridLayout
}

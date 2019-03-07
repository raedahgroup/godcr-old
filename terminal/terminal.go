package terminal

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
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

	creatWalletLayout := createWalletPage(tviewApp, ctx, walletMiddleware)
	return tviewApp.SetRoot(creatWalletLayout, false).Run()
}

func terminalLayout(tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {
	gridLayout := tview.NewGrid().
		SetRows(3, 0).
		SetColumns(25, 0, 1).
		SetGap(0, 2).
		SetBordersColor(helpers.DecredLightColor)

	gridLayout.SetBackgroundColor(tcell.ColorBlack)

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
		tviewApp.Draw()
		tviewApp.SetFocus(menuColumn)
	}

	menuColumn.AddItem("Overview", "", 'o', func() {
		changePageColumn(pages.BalancePage(walletMiddleware, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("History", "", 'h', func() {
		changePageColumn(pages.HistoryPage(walletMiddleware, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Send", "", 's', func() {
		changePageColumn(pages.SendPage(walletMiddleware, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Receive", "", 'r', func() {
		changePageColumn(pages.ReceivePage(walletMiddleware, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Staking", "", 'k', func() {
		changePageColumn(pages.StakingPage(walletMiddleware, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Accounts", "", 'a', nil)

	menuColumn.AddItem("Security", "", 'c', nil)

	menuColumn.AddItem("Settings", "", 't', nil)

	menuColumn.AddItem("Exit", "", 'q', func() {
		tviewApp.Stop()
	})

	header := primitives.NewCenterAlignedTextView(fmt.Sprintf("\n%s Terminal", strings.ToUpper(app.Name)))
	header.SetBackgroundColor(helpers.DecredColor)
	gridLayout.AddItem(header, 0, 0, 1, 3, 0, 0, false)

	menuColumn.SetShortcutColor(helpers.DecredLightColor)
	menuColumn.SetBorder(true).SetBorderColor(helpers.DecredLightColor)
	gridLayout.AddItem(menuColumn, 1, 0, 1, 1, 0, 0, true)

	menuColumn.SetCurrentItem(0)
	changePageColumn(pages.BalancePage(walletMiddleware, tviewApp.SetFocus, clearFocus))

	return gridLayout
}

func createWalletPage(tviewApp *tview.Application, ctx context.Context, walletMiddleware app.WalletMiddleware) tview.Primitive {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	
	layout.AddItem(tview.NewTextView().SetText("Create Wallet").SetTextAlign(tview.AlignCenter), 4, 1, false)
		
	// get seed and display to user
	seed, err := walletMiddleware.GenerateNewWalletSeed()
	if err != nil {
		return layout.AddItem(tview.NewTextView().SetText(err.Error()).SetTextAlign(tview.AlignCenter), 4, 1, false)
	}

	outputTextView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	outputMessage := func(text string) {
		layout.RemoveItem(outputTextView)
		layout.AddItem(outputTextView.SetText(text), 0, 1, true)
	}

	form := tview.NewForm()
	var password string
	form.AddPasswordField("PassPhrase", "", 20, '*', func(text string) {
		password = text
	})

	var confPassword string
	form.AddPasswordField("Confirm Password", "", 20, '*', func(text string) {
		confPassword = text
	})

	var hasStoredSeed bool
	form.AddCheckbox("I have stored wallet seed", false, func(checked bool) {
		hasStoredSeed = checked
	})

	form.AddButton("Create", func() {
		if len(password) == 0 {
			outputMessage("PassPhrase cannot empty")
			return
		}
		if password != confPassword {
			outputMessage("password does not match")
			return
		}
		if hasStoredSeed == false {
			outputMessage("please Store seed in a szfe location and check the box")
			return
		}

		err := CreateWallet(ctx, seed, password, walletMiddleware)
		if err != nil {
			layout.AddItem(tview.NewTextView().SetText(err.Error()).SetTextAlign(tview.AlignCenter), 4, 1, false)
			return
		}
	})

	form.AddButton("Quit", func() {
		tviewApp.Stop()
	})

	layout.AddItem(form, 10, 1, true)

	seedInfo := tview.NewTextView().SetText(`IMPORTANT: Keep the seed in a safe place as you will NOT be able to restore your wallet without it. Please keep in mind that anyone who has access to the seed can also restore your wallet thereby giving them access to all your funds, so it is imperative that you keep it in a secure location.`).SetRegions(true).SetWordWrap(true)
	seedInfo.SetBorder(true)
	layout.AddItem(seedInfo, 7, 1, false)

	seedView := tview.NewTextView().SetRegions(true).SetWordWrap(true).SetText(seed)
	seedView.SetBorder(true)
	layout.AddItem(seedView, 7, 1, false)

	layout.SetFullScreen(true).SetBorderPadding(3, 1, 6, 4)	

	layout.SetFullScreen(true)
	return layout
}
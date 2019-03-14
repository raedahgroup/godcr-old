package terminal

import (
	"context"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/pages"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func StartTerminalApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	tviewApp := tview.NewApplication()

	page := pageLoader(ctx, tviewApp, walletMiddleware)
	// `Run` blocks until app.Stop() is called before returning
	return tviewApp.SetRoot(page, true).Run()
}

var Report string
var Status walletcore.SyncStatus

func terminalLayout(ctx context.Context, tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {

	syncBlockchain(walletMiddleware)
	if Status == walletcore.SyncStatusError {
		msgOutput := fmt.Sprintf("Cannot display page. %s", Report)
		return tview.NewTextView().SetText(msgOutput).SetTextAlign(tview.AlignCenter)
	}
	if Status == walletcore.SyncStatusInProgress {
		msgOutput := fmt.Sprintf("%s. Refresh after a while to access this page", Report)
		return tview.NewTextView().SetText(msgOutput).SetTextAlign(tview.AlignCenter)
	}
	if Status == walletcore.SyncStatusSuccess {

	header := helpers.CenterAlignedTextView(fmt.Sprintf("\n%s Terminal", strings.ToUpper(app.Name)))
	header.SetBackgroundColor(helpers.DecredColor)
		//Creating the View for the Layout
	gridLayout := tview.NewGrid().SetRows(3, 0).SetColumns(30, 0)
	//Controls the display for the right side column
	var activePage tview.Primitive
	changePageColumn := func(page tview.Primitive) {
		gridLayout.RemoveItem(activePage)
		activePage = page
		gridLayout.AddItem(activePage, 1, 1, 1, 1, 0, 0, true)
	}
	var menuColumn *tview.List
	setFocus := tviewApp.SetFocus

	menuColumn := tview.NewList()
	clearFocus := func() {
		gridLayout.RemoveItem(activePage)
		tviewApp.SetFocus(menuColumn)
	}

	//Menu List of the Layout
	menuColumn.AddItem("Balance", "", 'b', func() {
		changePageColumn(pages.BalancePage(walletMiddleware, setFocus, clearFocus))
	})

	menuColumn.AddItem("Receive", "", 'r', func() {
		changePageColumn(pages.ReceivePage(walletMiddleware, setFocus, clearFocus))
	})

	menuColumn.AddItem("Send", "", 's', func() {
		changePageColumn(pages.SendPage(setFocus, clearFocus))
	})

	menuColumn.AddItem("History", "", 'h', func() {
		changePageColumn(pages.HistoryPage(walletMiddleware, setFocus, clearFocus))
	})

	menuColumn.AddItem("Stake Info", "", 'k', func() {
		changePageColumn(pages.StakeinfoPage(walletMiddleware, setFocus, clearFocus))
	})

	menuColumn.AddItem("Purchase Tickets", "", 't', func() {
		changePageColumn(pages.PurchaseTicketsPage(walletMiddleware, setFocus, clearFocus))
	})

	menuColumn.AddItem("Exit", "", 'q', func() {
		tviewApp.Stop()
	})

	menuColumn.SetCurrentItem(0)
	menuColumn.SetShortcutColor(helpers.DecredLightColor)
	menuColumn.SetBorder(true)
	menuColumn.SetBorderColor(helpers.DecredLightColor)
	// Layout for screens Header
	gridLayout.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	// Layout for screens with two column
	gridLayout.AddItem(menuColumn, 1, 0, 1, 1, 0, 0, true)
	changePageColumn(pages.BalancePage(walletMiddleware, setFocus, clearFocus))
	gridLayout.SetBorders

	return gridLayout

	} else {
		return tview.NewTextView().SetText("Cannot display page. Blockchain sync status cannot be determined").SetTextAlign(tview.AlignCenter)
	}
}

func pageLoader(ctx context.Context, tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {
	walletExists, err := openWalletIfExist(ctx, walletMiddleware)
	if err != nil {
		return tview.NewTextView().SetText(err.Error()).SetTextAlign(tview.AlignCenter)
	}

<<<<<<< HEAD
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
=======

	pages := tview.NewPages()
>>>>>>> added page primitive to for easy page navigation after wallet creation

	if walletExists {
		pages.AddPage("main", terminalLayout(ctx, tviewApp, walletMiddleware), true, true)
	} else {

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
		form.AddPasswordField("Passphrase", "", 20, '*', func(text string) {
			password = text
		})

		var confPassword string
		form.AddPasswordField("Confirm Passphrase", "", 20, '*', func(text string) {
			confPassword = text
		})

		var hasStoredSeed bool
		form.AddCheckbox("I have stored wallet seed", false, func(checked bool) {
			hasStoredSeed = checked
		})

		form.AddButton("Create", func() {
			if len(password) == 0 {
				outputMessage("Passphrase cannot empty")
				return
			}
			if password != confPassword {
				outputMessage("passphrase does not match")
				return
			}
			if hasStoredSeed == false {
				outputMessage("please Store seed in a szfe location and check the box")
				return
			}

			err := CreateWallet(tviewApp, seed, password, walletMiddleware)
			if err != nil {
				outputMessage(err.Error())
				return
			}

			pages.AddAndSwitchToPage("main", terminalLayout(ctx, tviewApp, walletMiddleware), true)
		})

		form.AddButton("Quit", func() {
			tviewApp.Stop()
		})

		layout.AddItem(form, 10, 1, true)

		seedInfo := tview.NewTextView().SetText(`Keep the seed in a safe place as you will NOT be able to restore your wallet without it. Please keep in mind that anyone who has access to the seed can also restore your wallet thereby giving them access to all your funds, so it is imperative that you keep it in a secure location.`).SetRegions(true).SetWordWrap(true)
		seedInfo.SetBorder(true).SetTitle("IMPORTANT:").SetTitleColor(helpers.WarnColor)
		layout.AddItem(seedInfo, 7, 1, false)

		seedView := tview.NewTextView().SetRegions(true).SetWordWrap(true).SetText(seed)
		seedView.SetBorder(true).SetTitle("Wallet Seed").SetTitleColor(helpers.PrimaryColor)
		layout.AddItem(seedView, 7, 1, false)

		layout.SetFullScreen(true).SetBorderPadding(3, 1, 6, 4)

		layout.SetFullScreen(true)

		pages.AddPage("create", layout, true, true)
	}

	return pages
}

func syncBlockchain(wallet app.WalletMiddleware) {

	err := wallet.SyncBlockChain(&app.BlockChainSyncListener{
		SyncStarted: func() {
			updateStatus("Blockchain sync started...", walletcore.SyncStatusInProgress)
		},
		SyncEnded: func(err error) {
			if err != nil {
				updateStatus(fmt.Sprintf("Blockchain sync completed with error: %s", err.Error()), walletcore.SyncStatusError)
			} else {
				updateStatus("Blockchain sync completed successfully", walletcore.SyncStatusSuccess)
			}
		},
		OnHeadersFetched: func(percentageProgress int64) {
			updateStatus(fmt.Sprintf("Blockchain sync in progress. Fetching headers (1/3): %d%%", percentageProgress), walletcore.SyncStatusInProgress)
		},
		OnDiscoveredAddress: func(_ string) {
			updateStatus("Blockchain sync in progress. Discovering addresses (2/3)", walletcore.SyncStatusInProgress)
		},
		OnRescanningBlocks: func(percentageProgress int64) {
			updateStatus(fmt.Sprintf("Blockchain sync in progress. Rescanning blocks (3/3): %d%%", percentageProgress), walletcore.SyncStatusInProgress)
		},
	}, false)

	if err != nil {
		updateStatus(fmt.Sprintf("Blockchain sync failed to start. %s", err.Error()), walletcore.SyncStatusError)
	}
}

func updateStatus(report string, status walletcore.SyncStatus) {
	Report = report
	Status = status
}

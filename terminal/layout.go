package terminal

import (
	"context"
	"fmt"
	"strings"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/pages"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/rivo/tview"
)

func terminalLayout(ctx context.Context, tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {
	syncBlockchain(walletMiddleware)
	if Status == walletcore.SyncStatusError {
		msgOutput := fmt.Sprintf(Report)
		helpers.CenterAlignedTextView(msgOutput)
		tviewApp.Stop()
	}
	if Status == walletcore.SyncStatusInProgress {
		msgOutput := fmt.Sprintf(Report)
		return helpers.CenterAlignedTextView(msgOutput)
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
		gridLayout.SetBordersColor(helpers.DecredLightColor)
		
		return gridLayout
		} else {
			return helpers.CenterAlignedTextView("Cannot display page. Blockchain sync status cannot be determined")
	}
}

var Report string
var Status walletcore.SyncStatus

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

package pages

import (
	"fmt"
	"strings"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func TerminalLayout(tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {
	header := primitives.NewCenterAlignedTextView(fmt.Sprintf("\n%s Terminal", strings.ToUpper(app.Name)))
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
		changePageColumn(BalancePage(walletMiddleware, setFocus, clearFocus))
	})

	menuColumn.AddItem("Receive", "", 'r', func() {
		changePageColumn(ReceivePage(walletMiddleware, setFocus, clearFocus))
	})

	menuColumn.AddItem("Send", "", 's', func() {
		changePageColumn(SendPage(walletMiddleware, setFocus, clearFocus))
	})

	menuColumn.AddItem("History", "", 'h', func() {
		changePageColumn(HistoryPage(walletMiddleware, setFocus, clearFocus))
	})

	menuColumn.AddItem("Staking", "", 'k', func() {
		changePageColumn(StakingPage(walletMiddleware, tviewApp.SetFocus, clearFocus))
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
	changePageColumn(BalancePage(walletMiddleware, setFocus, clearFocus))
	gridLayout.SetBordersColor(helpers.DecredLightColor)

	return gridLayout
}

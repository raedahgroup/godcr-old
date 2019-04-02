package pages

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func rootPage(tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {
	gridLayout := tview.NewGrid().
		SetRows(3, 0).
		SetColumns(25, 0, 1).
		SetGap(0, 2).
		SetBordersColor(helpers.DecredLightColor)

	gridLayout.SetBackgroundColor(tcell.ColorBlack)

	var activePage tview.Primitive

	// displayPage sets the currently active page to be displayed on the second column of the grid layout
	displayPage := func(page tview.Primitive) {
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
		displayPage(balancePage(walletMiddleware, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("History", "", 'h', func() {
		displayPage(historyPage(walletMiddleware, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Send", "", 's', func() {
		displayPage(sendPage(walletMiddleware, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Receive", "", 'r', func() {
		displayPage(receivePage(walletMiddleware, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Staking", "", 'k', func() {
		displayPage(stakingPage(walletMiddleware, tviewApp.SetFocus, clearFocus))
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
	displayPage(balancePage(walletMiddleware, tviewApp.SetFocus, clearFocus))

	return gridLayout
}

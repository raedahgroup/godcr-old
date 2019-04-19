package pages

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func rootPage(tviewApp *tview.Application, walletMiddleware app.WalletMiddleware) tview.Primitive {
	gridLayout := tview.NewGrid().
		SetRows(3, 1, 0, 1, 3).
		SetColumns(25, 2, 0, 2).
		SetBordersColor(helpers.DecredLightBlueColor)

	gridLayout.SetBackgroundColor(tcell.ColorBlack)

	var activePage tview.Primitive

	// displayPage sets the currently active page to be displayed on the second column of the grid layout
	displayPage := func(page tview.Primitive) {
		gridLayout.RemoveItem(activePage)
		activePage = page
		gridLayout.AddItem(activePage, 2, 2, 1, 1, 0, 0, true)
	}

	hintTextView := primitives.WordWrappedTextView("")
	hintTextView.SetTextColor(helpers.HintTextColor)

	menuColumn := tview.NewList()
	clearFocus := func() {
		gridLayout.RemoveItem(activePage)
		hintTextView.SetText("")
		tviewApp.Draw()
		tviewApp.SetFocus(menuColumn)
	}

	menuColumn.AddItem("Overview", "", 0, func() {
		displayPage(overviewPage(walletMiddleware, hintTextView, tviewApp, clearFocus))
	})

	menuColumn.AddItem("History", "", 0, func() {
		displayPage(historyPage(walletMiddleware, hintTextView, tviewApp, clearFocus))
	})

	menuColumn.AddItem("Send", "", 0, func() {
		displayPage(sendPage(walletMiddleware, hintTextView, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Receive", "", 0, func() {
		displayPage(receivePage(walletMiddleware, hintTextView, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Staking", "", 0, func() {
		displayPage(stakingPage(walletMiddleware, hintTextView, tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Accounts", "", 0, func() {
		displayPage(accountsPage(tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Security", "", 0, func() {
		displayPage(securityPage(tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Settings", "", 0, func() {
		displayPage(settingsPage(tviewApp.SetFocus, clearFocus))
	})

	menuColumn.AddItem("Exit", "", 0, func() {
		tviewApp.Stop()
	})

	netType := walletMiddleware.NetType()
	header := primitives.NewCenterAlignedTextView(fmt.Sprintf("\n %s %s\n", app.DisplayName, netType))
	header.SetBackgroundColor(helpers.DecredBlueColor)
	gridLayout.AddItem(header, 0, 0, 1, 4, 0, 0, false)

	menuColumn.SetShortcutColor(helpers.DecredLightBlueColor)
	menuColumn.SetBorder(true).SetBorderColor(helpers.DecredLightBlueColor)
	gridLayout.AddItem(menuColumn, 1, 0, 4, 1, 0, 0, true)

	gridLayout.AddItem(hintTextView, 4, 2, 1, 1, 0, 0, false)

	menuColumn.SetCurrentItem(0)

	displayPage(overviewPage(walletMiddleware, hintTextView, tviewApp, clearFocus))

	return gridLayout
}

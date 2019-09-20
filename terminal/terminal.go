package terminal

import (
	"fmt"
	"os"

	"github.com/decred/slog"
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/pages"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

type terminalUI struct {
	app               *tview.Application
	rootGridLayout    *tview.Grid
	navMenu           *primitives.List
	pageContentHolder *tview.Flex
	activePageContent tview.Primitive
	hintTextView      *primitives.TextView
	log               slog.Logger
	dcrlw             *dcrlibwallet.LibWallet
}

func LaunchUserInterface(appDisplayName, defaultAppDataDir, netType string) {
	logger, err := dcrlibwallet.RegisterLogger("TUI")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Launch error - cannot register logger: %v", err)
		return
	}

	tui := &terminalUI{
		app: tview.NewApplication(),
		log: logger,
	}

	tui.dcrlw, err = dcrlibwallet.NewLibWallet(defaultAppDataDir, "", netType)
	if err != nil {
		tui.log.Errorf("Initialization error: %v", err)
		return
	}

	walletExists, err := tui.dcrlw.WalletExists()
	if err != nil {
		tui.log.Errorf("Error checking if wallet db exists: %v", err)
		return
	}

	if !walletExists {
		// todo show create wallet page
		tui.log.Infof("Wallet does not exist in app directory. Need to create one.")
		return
	}

	var pubPass []byte = nil // todo is this initial value setting necessary?
	if tui.dcrlw.ReadBoolConfigValueForKey(dcrlibwallet.IsStartupSecuritySetConfigKey) {
		// prompt user for public passphrase and assign to `pubPass`
	}

	err = tui.dcrlw.OpenWallet(pubPass)
	if err != nil {
		tui.log.Errorf("Error opening wallet db: %v", err)
		return
	}

	err = tui.dcrlw.SpvSync("") // todo dcrlibwallet should ideally read this parameter from config
	if err != nil {
		tui.log.Errorf("Spv sync attempt failed: %v", err)
		return
	}

	tui.prepareMainWindow(appDisplayName)
	tui.app.SetRoot(tui.rootGridLayout, true)

	// app is ready, pass necessary variables to pages pkg and display first page
	pages.Setup(tui.app, tui.log, tui.dcrlw, tui.hintTextView, tui.clearPageContent)
	tui.navMenu.SetCurrentItem(0)

	// turn off all logging at this point
	dcrlibwallet.SetLogLevels("off")

	if err = tui.app.Run(); err != nil {
		tui.log.Errorf("App exited due to error: %v", err)
	}
}

func (tui *terminalUI) prepareMainWindow(appDisplayName string) {
	/*
		todo correct this drawing
		| Godcr    |  |       <title>       |  |    |
		|          |  |                     |  |    |
		| nav menu |  | <main page content> |  |    |
		| nav menu |  | <main page content> |  |    |
		| nav menu |  | <main page content> |  |    |
		| nav menu |  | <main page content> |  |    |
		|          |  |       <footer>      |  |    |
	*/
	tui.rootGridLayout = tview.NewGrid()
	tui.rootGridLayout.SetRows(3, 1, 0, 1, 2)  // nav menu, space, main content (max width), space, space(?)
	tui.rootGridLayout.SetColumns(20, 2, 0, 2) // title, space, main content (max height), space
	tui.rootGridLayout.SetBackgroundColor(tcell.ColorBlack)

	// row 0, col 0, span 1 row and all (4) columns
	headerText := fmt.Sprintf("\n %s %s\n", appDisplayName, tui.dcrlw.NetType())
	header := primitives.NewCenterAlignedTextView(headerText).SetBackgroundColor(helpers.DecredBlueColor)
	tui.rootGridLayout.AddItem(header, 0, 0, 1, 4, 0, 0, false)

	// row 1, col 0, span 4 rows and 1 column
	tui.prepareNavigationMenu()
	tui.rootGridLayout.AddItem(tui.navMenu, 1, 0, 4, 1, 0, 0, true)

	// row 1, col 1, span 3 columns and 3 rows HUH???
	tui.pageContentHolder = tview.NewFlex().SetDirection(tview.FlexRow)
	tui.pageContentHolder.SetBorderColor(helpers.DecredLightBlueColor)
	tui.pageContentHolder.SetBorderPadding(0, 0, 1, 1)
	tui.rootGridLayout.AddItem(tui.pageContentHolder, 1, 1, 3, 3, 0, 0, true)

	// row 4, col 1, space 3 columns
	tui.hintTextView = primitives.WordWrappedTextView("")
	tui.hintTextView.SetTextColor(helpers.HintTextColor)
	tui.rootGridLayout.AddItem(tui.hintTextView, 4, 1, 1, 3, 0, 0, false)
}

func (tui *terminalUI) prepareNavigationMenu() {
	tui.navMenu = primitives.NewList()
	tui.navMenu.SetBorderColor(helpers.DecredLightBlueColor)

	for _, page := range pages.All() {
		tui.navMenu.AddItem(page.Name, "", page.Shortcut, func() {
			tui.clearPageContent()
			tui.removeNavMenuFocus()
			tui.setPageContent(page.Content())
		})
	}

	// todo escape button listener, should exit app instead
	tui.navMenu.SetDoneFunc(func() {
		tui.clearPageContent()
		tui.focusNavMenu()
		tui.app.Draw()
	})
}

func (tui *terminalUI) focusNavMenu() {
	tui.navMenu.SetBorder(true)
	tui.navMenu.ShowShortcut(true)
	tui.navMenu.SetMainTextColor(tcell.ColorWhite)
	tui.navMenu.SetBorderPadding(0, 0, 1, 0)
	tui.app.SetFocus(tui.navMenu)
}

func (tui *terminalUI) removeNavMenuFocus() {
	tui.navMenu.SetBorder(false)
	tui.navMenu.ShowShortcut(false)
	tui.navMenu.SetMainTextColor(helpers.HintTextColor)
	tui.navMenu.SetBorderPadding(1, 0, 1, 0)
}

func (tui *terminalUI) clearPageContent() {
	tui.hintTextView.SetText("")
	tui.pageContentHolder.RemoveItem(tui.activePageContent)
	tui.pageContentHolder.SetBorder(false)
}

func (tui *terminalUI) setPageContent(pageContent tview.Primitive) {
	tui.activePageContent = pageContent
	tui.pageContentHolder.AddItem(tui.activePageContent, 0, 1, true)
	tui.pageContentHolder.SetBorder(true)
}

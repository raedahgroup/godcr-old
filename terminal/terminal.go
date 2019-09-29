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

	var pubPass []byte
	if tui.dcrlw.ReadBoolConfigValueForKey(dcrlibwallet.IsStartupSecuritySetConfigKey) {
		// prompt user for public passphrase and assign to `pubPass`
	}

	err = tui.dcrlw.OpenWallet(pubPass)
	if err != nil {
		tui.log.Errorf("Error opening wallet db: %v", err)
		return
	}

	err = tui.dcrlw.SpvSync("")
	if err != nil {
		tui.log.Errorf("Spv sync attempt failed: %v", err)
		return
	}

	tui.prepareMainWindow(appDisplayName)
	tui.app.SetRoot(tui.rootGridLayout, true)

	// app is ready, pass necessary variables to pages pkg and display first page
	pages.Setup(tui.app, tui.log, tui.dcrlw, tui.hintTextView, tui.clearPageContent)
	firstPageContent := pages.All()[0].Content()
	tui.removeNavMenuFocus()
	tui.setPageContent(firstPageContent)

	// turn off all logging at this point
	dcrlibwallet.SetLogLevels("off")

	err = tui.app.Run()

	// app has exited, revert to using default log level
	userConfiguredLogLevel := tui.dcrlw.ReadStringConfigValueForKey(dcrlibwallet.LogLevelConfigKey)
	if userConfiguredLogLevel == "" {
		userConfiguredLogLevel = "info"
	}
	dcrlibwallet.SetLogLevels(userConfiguredLogLevel)

	if err != nil {
		tui.log.Errorf("App exited due to error: %v\nShutting down wallet...", err)
	} else {
		tui.log.Infof("App exited. Shutting down wallet...")
	}

	tui.dcrlw.Shutdown()
}

func (tui *terminalUI) prepareMainWindow(appDisplayName string) {
	/*
		------------------------------------
		|          GoDCR {network}         | -> row 0
		------------------------------------
		| nav menu |  <main page content>  | -> row 1
		           -------------------------
		| nav menu |  <current page hint>  | -> row 2
		------------------------------------
		| ^^ col 0 |       ^^ col 1        |
	*/
	tui.rootGridLayout = tview.NewGrid()
	tui.rootGridLayout.SetRows(3, 0, 2)
	tui.rootGridLayout.SetColumns(20, 0)
	tui.rootGridLayout.SetBackgroundColor(tcell.ColorBlack)

	// app name -> row 0 (first row), col 0 to col 1 (span 2 columns)
	headerText := fmt.Sprintf("\n %s %s\n", appDisplayName, tui.dcrlw.NetType())
	header := primitives.NewCenterAlignedTextView(headerText) /*.SetBackgroundColor(helpers.DecredBlueColor)*/
	tui.rootGridLayout.AddItem(header, 0, 0, 1, 2, 0, 0, false)

	// nav menu -> row 1 to row 2 (span 2 rows), col 0 (first column)
	tui.prepareNavigationMenu()
	tui.rootGridLayout.AddItem(tui.navMenu, 1, 0, 3, 1, 0, 0, true)

	// row 1 (middle row), col 1 (second column)
	tui.pageContentHolder = tview.NewFlex().SetDirection(tview.FlexRow)
	tui.pageContentHolder.SetBorderColor(helpers.DecredLightBlueColor)
	tui.pageContentHolder.SetBorderPadding(0, 0, 1, 1)
	tui.rootGridLayout.AddItem(tui.pageContentHolder, 1, 1, 3, 1, 0, 0, true)

	// row 2 (bottom row), col 1 (second column)
	tui.hintTextView = primitives.WordWrappedTextView("")
	tui.hintTextView.SetTextColor(helpers.HintTextColor)
	tui.rootGridLayout.AddItem(tui.hintTextView, 2, 1, 1, 1, 0, 0, false)
}

func (tui *terminalUI) prepareNavigationMenu() {
	tui.navMenu = primitives.NewList()
	tui.navMenu.SetBorderColor(helpers.DecredLightBlueColor)

	for _, page := range pages.All() {
		tui.navMenu.AddItem(page.Name, "", page.Shortcut, tui.makePageContentLoaderFn(page.Content))
	}

	// escape button from main menu should display exit page prompt
	tui.navMenu.SetDoneFunc(tui.makePageContentLoaderFn(pages.ExitPage))
}

// makePageContentLoaderFn returns a func that, when triggered,
// creates the display for a page and renders it.
func (tui *terminalUI) makePageContentLoaderFn(pageContentFn func() tview.Primitive) func() {
	return func() {
		tui.clearPageContent()
		tui.removeNavMenuFocus()
		tui.setPageContent(pageContentFn())
	}
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
	tui.focusNavMenu()
}

func (tui *terminalUI) setPageContent(pageContent tview.Primitive) {
	tui.activePageContent = pageContent
	tui.pageContentHolder.AddItem(tui.activePageContent, 0, 1, true)
	tui.pageContentHolder.SetBorder(true)
}

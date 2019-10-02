package nuklear

import (
	"fmt"
	"image"
	"os"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/nuklear/nuklog"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

const (
	navWidth            = 200
	defaultWindowWidth  = 800
	defaultWindowHeight = 600
)

type nuklearApp struct {
	appDisplayName string
	wallet         *dcrlibwallet.LibWallet
	syncer         *Syncer
	navPages       map[string]navPageHandler
	currentPage    string
	nextPage       string
	pageChanged    bool
}

func LaunchUserInterface(appDisplayName, appDataDir, netType string) {
	logger, err := dcrlibwallet.RegisterLogger("NUKL")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Launch error - cannot register logger: %v", err)
		return
	}

	nuklog.UseLogger(logger)

	app := &nuklearApp{}

	app.wallet, err = dcrlibwallet.NewLibWallet(appDataDir, "", netType)
	if err != nil {
		// todo display pre-launch error on UI
		nuklog.Log.Errorf("Initialization error: %v", err)
		return
	}

	walletExists, err := app.wallet.WalletExists()
	if err != nil {
		// todo display pre-launch error on UI
		nuklog.Log.Errorf("Error checking if wallet db exists: %v", err)
		return
	}

	if !walletExists {
		// todo show create wallet page
		nuklog.Log.Infof("Wallet does not exist in app directory. Need to create one.")
		return
	}

	var pubPass []byte
	if app.wallet.ReadBoolConfigValueForKey(dcrlibwallet.IsStartupSecuritySetConfigKey) {
		// prompt user for public passphrase and assign to `pubPass`
	}

	err = app.wallet.OpenWallet(pubPass)
	if err != nil {
		// todo display pre-launch error on UI
		nuklog.Log.Errorf("Error opening wallet db: %v", err)
		return
	}

	err = app.wallet.SpvSync("")
	if err != nil {
		// todo display pre-launch error on UI
		nuklog.Log.Errorf("Spv sync attempt failed: %v", err)
		return
	}

	// initialize fonts for later use
	err = styles.InitFonts()
	if err != nil {
		// todo display pre-launch error on UI
		nuklog.Log.Errorf("Error initializing app fonts: %v", err)
		return
	}

	// initialize master window and set style
	windowSize := image.Point{X: defaultWindowWidth, Y: defaultWindowHeight}
	masterWindow := nucular.NewMasterWindowSize(nucular.WindowNoScrollbar, appDisplayName, windowSize, app.render)
	masterWindow.SetStyle(styles.MasterWindowStyle())

	// register nav page handlers
	navPages := getNavPages()
	app.navPages = make(map[string]navPageHandler, len(navPages))
	for _, page := range navPages {
		app.navPages[page.name] = page.handler
	}

	app.syncer = NewSyncer(app.wallet, masterWindow.Changed)
	app.wallet.AddSyncProgressListener(app.syncer, app.appDisplayName)

	app.currentPage = "overview"
	app.pageChanged = true

	// draw master window
	masterWindow.Main()
}

func (app *nuklearApp) render(window *nucular.Window) {
	if _, isNavPage := app.navPages[app.currentPage]; isNavPage {
		app.renderNavPage(window)
	} else {
		errorMessage := fmt.Sprintf("Page not properly set up: %s", app.currentPage)
		nuklog.Log.Errorf(errorMessage)

		w := &widgets.Window{window}
		w.DisplayMessage(errorMessage, styles.DecredOrangeColor)
	}
}

func (app *nuklearApp) renderNavPage(window *nucular.Window) {
	// this creates the space on the window that will hold 2 widgets
	// the navigation section on the window and the main page content
	entireWindow := window.Row(0).SpaceBegin(2)

	app.renderNavWindow(window, entireWindow.H)
	app.renderPageContentWindow(window, entireWindow.W, entireWindow.H)
}

func (app *nuklearApp) renderNavWindow(window *nucular.Window, maxHeight int) {
	navSectionRect := rect.Rect{
		X: 0,
		Y: 0,
		W: navWidth,
		H: maxHeight,
	}
	window.LayoutSpacePushScaled(navSectionRect)

	// set style
	styles.SetNavStyle(window.Master())

	// create window and draw nav menu
	widgets.NoScrollGroupWindow("nav-group-window", window, func(navGroupWindow *widgets.Window) {
		navGroupWindow.AddHorizontalSpace(10)
		navGroupWindow.AddColoredLabel(fmt.Sprintf("%s %s", app.appDisplayName, app.wallet.NetType()),
			styles.DecredLightBlueColor, widgets.CenterAlign)
		navGroupWindow.AddHorizontalSpace(10)

		for _, page := range getNavPages() {
			if app.currentPage == page.name {
				navGroupWindow.AddCurrentNavButton(page.label, func() {
					app.changePage(window, page.name)
				})
			} else {
				navGroupWindow.AddBigButton(page.label, func() {
					app.changePage(window, page.name)
				})
			}
		}

		// add exit button
		navGroupWindow.AddBigButton("Exit", func() {
			go navGroupWindow.Master().Close()
		})
	})
}

func (app *nuklearApp) renderPageContentWindow(window *nucular.Window, maxWidth, maxHeight int) {
	pageSectionRect := rect.Rect{
		X: navWidth,
		Y: 0,
		W: maxWidth - navWidth,
		H: maxHeight,
	}

	// set style
	styles.SetPageStyle(window.Master())
	window.LayoutSpacePushScaled(pageSectionRect)

	if app.wallet.IsSyncing() {
		app.syncer.Render(window)
	} else {
		handler := app.navPages[app.currentPage]
		// ensure that the handler's BeforeRender function is called only once per page call
		// as it initializes page variables
		if app.pageChanged {
			handler.BeforeRender(app.wallet, window.Master().Changed)
			app.pageChanged = false
		}
		handler.Render(window)
	}
}

func (app *nuklearApp) changePage(window *nucular.Window, newPage string) {
	app.nextPage = newPage
	app.currentPage = newPage
	app.pageChanged = true
	window.Master().Changed()
}

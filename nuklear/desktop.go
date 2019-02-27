package nuklear

import (
	"context"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/rect"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

const (
	navWidth  = 260
	homePage  = "balance"
	startPage = "sync"

	contentPaneXOffset      = 45
	contentPaneWidthPadding = 55
)

type Desktop struct {
	masterWindow    nucular.MasterWindow
	wallet          app.WalletMiddleware
	currentPage     string
	pageChanged     bool
	navPages        map[string]NavPageHandler
	standalonePages map[string]StandalonePageHandler
}

func LaunchApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	desktop := &Desktop{
		walletMiddleware: walletMiddleware,
		pageChanged:      true,
		currentPage:      homePage,
	}

	// initialize master window and set style
	window := nucular.NewMasterWindow(nucular.WindowNoScrollbar, app.Name, desktop.render)
	window.SetStyle(helpers.GetStyle())
	desktop.masterWindow = window

	// initialize fonts for later use
	err := helpers.InitFonts()
	if err != nil {
		return nil
	}

	// register nav page handlers
	navPages := getNavPagesData()
	desktop.navPages = make(map[string]NavPageHandler, len(navPages))
	for _, page := range navPages {
		desktop.navPages[page.name] = page.handler
	}

	// register standalone page handlers
	standalonePages := getStandalonePageData()
	desktop.standalonePages = make(map[string]StandalonePageHandler, len(standalonePages))
	for _, page := range standalonePages {
		desktop.standalonePages[page.name] = page.handler
	}

	// open wallet and start blockchain syncing in background
	walletExists, err := openWalletIfExist(ctx, walletMiddleware)
	if err != nil {
		return err
	}

	if !walletExists {
		// todo add ui to create wallet
		err = fmt.Errorf("No wallet found. Use 'godcr create' to create a wallet before launching the desktop app")
		nuklog.LogInfo(err.Error())
		return err
	}

	desktop.currentPage = startPage

	// draw window
	desktop.masterWindow.Main()
	return nil
}

func (desktop *Desktop) render(window *nucular.Window) {
	if handler, ok := desktop.standalonePages[desktop.currentPage]; ok {
		desktop.renderStandalonePage(window, handler)
		return
	}

	handler := desktop.navPages[desktop.currentPage]
	desktop.renderNavPage(window, handler)
}

func (desktop *Desktop) renderNavPage(window *nucular.Window, handler NavPageHandler) {
	area := window.Row(0).SpaceBegin(2)

	// create nav pane
	navRect := rect.Rect{
		X: 0,
		Y: 0,
		W: navWidth,
		H: area.H,
	}
	window.LayoutSpacePushScaled(navRect)

	// render nav
	helpers.SetNavStyle(desktop.masterWindow)
	if navWindow := helpers.NewWindow("Navigation Group", window, 0); navWindow != nil {
		navWindow.Row(40).Dynamic(1)
		for _, page := range getNavPagesData() {
			if navWindow.Button(label.TA(page.label, "LC"), false) {
				desktop.changePage(page.name)
			}
		}
		navWindow.End()
	}

	// create content pane
	contentRect := rect.Rect{
		X: navWidth - contentPaneXOffset,
		Y: 0,
		W: (area.W + contentPaneWidthPadding) - navWidth,
		H: area.H,
	}

	helpers.SetContentArea(contentRect)

	// style content area
	helpers.SetPageStyle(desktop.masterWindow)

	window.LayoutSpacePushScaled(contentRect)
	if desktop.currentPage == "" { // ideally, this should only be false once in the lifetime of an instance
		desktop.changePage(homePage)
		return
	}

	// ensure that the handler's BeforeRender function is called only once per page call
	// as it initializes page variables
	if desktop.pageChanged {
		handler.BeforeRender()
		desktop.pageChanged = false
	}

	handler.Render(window, desktop.wallet)
}

func (desktop *Desktop) renderStandalonePage(window *nucular.Window, handler StandalonePageHandler) {
	window.Row(0).SpaceBeginRatio(1)
	window.LayoutSpacePushRatio(0.1, 0.05, 0.9, 0.8)

	if desktop.pageChanged {
		handler.BeforeRender()
		desktop.pageChanged = false
	}

	helpers.SetStandaloneWindowStyle(window.Master())
	handler.Render(window, desktop.wallet, desktop.changePage)
}

func (desktop *Desktop) changePage(page string) {
	desktop.currentPage = page
	desktop.pageChanged = true
	desktop.masterWindow.Changed()
}

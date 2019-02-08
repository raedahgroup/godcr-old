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
	navWidth = 260
	homePage = "balance"

	contentPaneXOffset      = 45
	contentPaneWidthPadding = 55
)

type Desktop struct {
	masterWindow     nucular.MasterWindow
	walletMiddleware app.WalletMiddleware
	currentPage      string
	pageChanged      bool
	pages            map[string]page
}

func LaunchApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	desktop := &Desktop{
		walletMiddleware: walletMiddleware,
		pageChanged:      true,
		currentPage:      homePage,
	}

	// register pages
	pages := getPages()
	desktop.pages = make(map[string]page, len(pages))
	for _, page := range pages {
		desktop.pages[page.name] = page
	}

	// open wallet and start blockchain syncing in background
	walletExists, err := openWalletIfExist(ctx, walletMiddleware)
	if err != nil {
		return err
	}

	if !walletExists {
		desktop.currentPage = "createwallet"
	}

	// initialize master window and set style
	window := nucular.NewMasterWindow(nucular.WindowNoScrollbar, app.Name, desktop.render)
	window.SetStyle(helpers.GetStyle())
	desktop.masterWindow = window

	// initialize fonts for later use
	err = helpers.InitFonts()
	if err != nil {
		return err
	}

	// todo run sync and show progress

	// draw window
	desktop.masterWindow.Main()
	return nil
}

func (desktop *Desktop) render(window *nucular.Window) {
	page := desktop.pages[desktop.currentPage]
	if page.standalone {
		desktop.renderStandalonePage(window, page.handler)
		return
	}

	desktop.renderNavPage(window, page.handler)
}

func (desktop *Desktop) renderStandalonePage(window *nucular.Window, handler Handler) {
	window.Row(0).SpaceBeginRatio(1)
	window.LayoutSpacePushRatio(0.1, 0.05, 0.9, 0.8)

	if desktop.pageChanged {
		handler.BeforeRender()
		desktop.pageChanged = false
	}

	helpers.SetStandaloneWindowStyle(window.Master())
	handler.Render(window, desktop.walletMiddleware, desktop.changePage)
}

func (desktop *Desktop) renderNavPage(window *nucular.Window, handler Handler) {
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
		for _, page := range getPages() {
			if !page.standalone {
				if navWindow.Button(label.TA(page.navLabel, "LC"), false) {
					desktop.changePage(page.name)
				}
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

	// ensure that the handler's BeforeRender function is called only once per page call
	// as it initializes page variables
	if desktop.pageChanged {
		handler.BeforeRender()
		desktop.pageChanged = false
	}

	handler.Render(window, desktop.walletMiddleware, desktop.changePage)
}

func (desktop *Desktop) changePage(page string) {
	desktop.currentPage = page
	desktop.pageChanged = true
	desktop.masterWindow.Changed()
}

package nuklear

import (
	"context"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/rect"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/nuklear/handlers"
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
	handlers         map[string]Handler
}

func LaunchApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	desktop := &Desktop{
		walletMiddleware: walletMiddleware,
	}

	// register handlers
	handlers := getHandlers()
	desktop.handlers = make(map[string]Handler, len(handlers))
	for _, handler := range handlers {
		desktop.handlers[handler.name] = handler.handler
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
	window := nucular.NewMasterWindow(nucular.WindowNoScrollbar, app.Name, desktop.renderWithNav)
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

func (desktop *Desktop) renderWithoutNav(window *nucular.Window) {
	window.Row(0).SpaceBeginRatio(1)
	window.LayoutSpacePushRatio(0.1, 0.05, 0.9, 0.8)

	handler := &handlers.CreateWalletHandler{} // for now, only this handler uses the renderWithoutNav function
	// ensure that the handler's BeforeRender function is called only once per page call
	// as it initializes page variables
	if desktop.pageChanged {
		handler.BeforeRender()
		desktop.pageChanged = false
	}

	helpers.SetNoNavWindowStyle(window.Master())
	handler.Render(window, desktop.walletMiddleware)
}

func (desktop *Desktop) renderWithNav(w *nucular.Window) {
	if desktop.currentPage == "createwallet" {
		desktop.renderWithoutNav(w)
		return
	}

	area := w.Row(0).SpaceBegin(2)

	// create nav pane
	navRect := rect.Rect{
		X: 0,
		Y: 0,
		W: navWidth,
		H: area.H,
	}
	w.LayoutSpacePushScaled(navRect)

	// render nav
	helpers.SetNavStyle(desktop.masterWindow)
	if navWindow := helpers.NewWindow("Navigation Group", w, 0); navWindow != nil {
		navWindow.Row(40).Dynamic(1)
		for _, handler := range getHandlers() {
			if navWindow.Button(label.TA(handler.navLabel, "LC"), false) {
				desktop.changePage(handler.name)
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

	w.LayoutSpacePushScaled(contentRect)
	if desktop.currentPage == "" { // ideally, this should only be false once in the lifetime of an instance
		desktop.changePage(homePage)
		return
	}

	handler := desktop.handlers[desktop.currentPage]

	// ensure that the handler's BeforeRender function is called only once per page call
	// as it initializes page variables
	if desktop.pageChanged {
		handler.BeforeRender()
		desktop.pageChanged = false
	}

	handler.Render(w, desktop.walletMiddleware)
}

func (desktop *Desktop) changePage(page string) {
	desktop.currentPage = page
	desktop.pageChanged = true
	desktop.masterWindow.Changed()
}

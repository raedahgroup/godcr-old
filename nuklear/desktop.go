package nuklear

import (
	"context"
	"fmt"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/rect"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

const (
	navWidth = 260
	homePage = "balance"
)

type Desktop struct {
	masterWindow nucular.MasterWindow
	wallet       walletcore.Wallet
	currentPage  string
	pageChanged  bool
	handlers     map[string]Handler
}

func LaunchApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	desktop := &Desktop{
		wallet: walletMiddleware,
	}

	// initialize master window and set style
	window := nucular.NewMasterWindow(nucular.WindowNoScrollbar, app.Name, desktop.render)
	window.SetStyle(helpers.GetStyle())
	desktop.masterWindow = window

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
		// todo add ui to create wallet
		err = fmt.Errorf("No wallet found. Use 'godcr create' to create a wallet before launching the desktop app")
		fmt.Println(err.Error())
		return err
	}

	// todo run sync and show progress

	// draw window
	desktop.masterWindow.Main()
	return nil
}

func (desktop *Desktop) render(w *nucular.Window) {
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
	if contentWindow := helpers.NewWindow("Navigation Group", w, 0); contentWindow != nil {
		contentWindow.Row(40).Dynamic(1)
		for _, handler := range getHandlers() {
			if contentWindow.Button(label.TA(handler.navLabel, "LC"), false) {
				desktop.changePage(handler.name)
			}
		}
		contentWindow.GroupEnd()
	}

	// create content pane
	contentRect := rect.Rect{
		X: navWidth - 45,
		Y: 0,
		W: (area.W + 55) - navWidth,
		H: area.H,
	}

	helpers.SetContentArea(contentRect)

	// style content area
	helpers.SetPageStyle(desktop.masterWindow)

	w.LayoutSpacePushScaled(contentRect)
	handler := desktop.handlers[desktop.currentPage]
	if desktop.currentPage == "" { // ideally, this should only be false once in the lifetime of an instance
		desktop.changePage(homePage)
		return
	}

	// ensure that the handler's BeforeRender function is called only once per page call
	if desktop.pageChanged {
		handler.BeforeRender()
		desktop.pageChanged = false
	}

	handler.Render(w, desktop.wallet)
}

func (desktop *Desktop) changePage(page string) {
	desktop.currentPage = page
	desktop.pageChanged = true
	desktop.masterWindow.Changed()
}

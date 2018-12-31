package desktop

import (
	"context"
	"fmt"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/rect"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
)

type pageHandler func(*nucular.Window)

type Desktop struct {
	window       nucular.MasterWindow
	currentPage  string
	wallet       walletcore.Wallet
	pageHandlers map[string]pageHandler
}

const (
	navWidth = 260
	homePage = "balance"
)

var (
	contentArea rect.Rect
)

func StartDesktopApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	d := &Desktop{
		wallet:       walletMiddleware,
		pageHandlers: make(map[string]pageHandler),
	}

	window := nucular.NewMasterWindow(nucular.WindowNoScrollbar, app.Name(), d.updateFn)
	window.SetStyle(getStyle())
	d.window = window

	d.registerHandlers()
	d.currentPage = homePage

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
	d.window.Main()
	return nil
}

func (d *Desktop) updateFn(w *nucular.Window) {
	area := w.Row(0).SpaceBegin(2)

	d.createNavPane(w, area.H)
	d.createContentPane(w, area.H, area.W)
}

func (d *Desktop) registerHandlers() {
	d.pageHandlers["balance"] = d.BalanceHandler
	d.pageHandlers["receive"] = d.ReceiveHandler
	d.pageHandlers["send"] = d.SendHandler
	d.pageHandlers["transactions"] = d.TransactionsHandler

	d.pageHandlers["selectutxos"] = d.selectUTXOSHandler
	d.pageHandlers["generateaddress"] = d.generateAddressHandler
}

func (d *Desktop) changePage(page string) {
	d.currentPage = page
	d.window.Changed()
}

func (d *Desktop) gotoPage(page string) {
	resetVars()
	d.changePage(page)
}

func (d *Desktop) gotoSubpage(page string) {
	d.changePage(page)
}

func (d *Desktop) createNavPane(w *nucular.Window, height int) {
	rect := rect.Rect{
		X: 0,
		Y: 0,
		W: navWidth,
		H: height,
	}

	w.LayoutSpacePushScaled(rect)
	// style navigation pane
	setNavStyle(d.window)
	if sw := w.GroupBegin("Navigation Group", 0); sw != nil {
		sw.Row(40).Dynamic(1)
		if sw.Button(label.TA("Balance", "LC"), false) {
			d.gotoPage("balance")
		}
		if sw.Button(label.TA("Send (WIP)", "LC"), false) {
			d.gotoPage("send")
		}
		if sw.Button(label.TA("Receive", "LC"), false) {
			d.gotoPage("receive")
		}
		if sw.Button(label.TA("Transactions", "LC"), false) {
			d.gotoPage("transactions")
		}
		sw.GroupEnd()
	}
}

func (d *Desktop) createContentPane(w *nucular.Window, height, width int) {
	rect := rect.Rect{
		X: navWidth - 45,
		Y: 0,
		W: (width + 55) - navWidth,
		H: height,
	}

	// style content area
	d.setPageStyle()

	contentArea = rect
	w.LayoutSpacePushScaled(rect)

	if pageHandler, ok := d.pageHandlers[d.currentPage]; ok {
		pageHandler(w)
	} else {
		d.pageHandlers[homePage](w)
	}
}

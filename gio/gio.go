package gio

import (
	"fmt"
	"image"
	"log"
	"os"

	gioapp "gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/io/system"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/gio/giolog"
	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
)

type (
	desktop struct {
		window         *gioapp.Window
		displayName    string
		pages          []page
		currentPage    string
		pageChanged    bool
		theme          *helper.Theme
		appDisplayName string
		multiWallet    *dcrlibwallet.MultiWallet
		syncer         *Syncer
	}
)

const (
	windowWidth  = 450
	windowHeight = 350

	navSectionWidth = 120
)

func LaunchUserInterface(appDisplayName, appDataDir, netType string) {
	logger, err := dcrlibwallet.RegisterLogger("GIOL")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Launch error - cannot register logger: %v", err)
		return
	}

	giolog.UseLogger(logger)

	theme := helper.NewTheme()
	app := &desktop{
		theme:       theme,
		currentPage: "overview",
	}

	app.multiWallet, err = dcrlibwallet.NewMultiWallet(appDataDir, "", netType)
	if err != nil {
		// todo display pre-launch error on UI
		giolog.Log.Errorf("Initialization error: %v", err)
		return
	}

	walletCount := app.multiWallet.LoadedWalletsCount()
	if walletCount == 0 {
		// todo show createand restore wallet page
		giolog.Log.Infof("Wallet does not exist in app directory. Need to create one.")
		return
	}

	var pubPass []byte
	if app.multiWallet.ReadBoolConfigValueForKey(dcrlibwallet.IsStartupSecuritySetConfigKey, true) {
		// prompt user for public passphrase and assign to `pubPass`
	}

	err = app.multiWallet.OpenWallets(pubPass)
	if err != nil {
		// todo display pre-launch error on UI
		giolog.Log.Errorf("Error opening wallet db: %v", err)
		return
	}

	err = app.multiWallet.SpvSync()
	if err != nil {
		// todo display pre-launch error on UI
		giolog.Log.Errorf("Spv sync attempt failed: %v", err)
		return
	}

	app.syncer = NewSyncer(theme, app.multiWallet, app.refreshWindow)
	app.multiWallet.AddSyncProgressListener(app.syncer, app.appDisplayName)

	app.prepareHandlers()
	go func() {
		app.window = gioapp.NewWindow(
			gioapp.Size(unit.Dp(windowWidth), unit.Dp(windowHeight)),
			gioapp.Title(app.displayName),
		)

		if err := app.renderLoop(); err != nil {
			log.Fatal(err)
		}
	}()

	// run app
	gioapp.Main()
}

func (d *desktop) prepareHandlers() {
	pages := getPages()
	d.pages = make([]page, len(pages))

	for index, page := range pages {
		d.pages[index] = page

		if index == 0 {
			d.changePage(page.name)
		}
	}
}

func (d *desktop) changePage(pageName string) {
	if d.currentPage == pageName {
		return
	}

	for _, page := range d.pages {
		if page.name == pageName {
			d.currentPage = page.name
			d.pageChanged = true
			break
		}
	}

}

func (d *desktop) renderLoop() error {
	ctx := &layout.Context{
		Queue: d.window.Queue(),
	}

	for {
		e := <-d.window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			ctx.Reset(e.Config, e.Size)
			d.render(ctx)
			e.Frame(ctx.Ops)
		}
	}
}

func (d *desktop) render(ctx *layout.Context) {
	var page page
	for i := range d.pages {
		if d.pages[i].name == d.currentPage {
			page = d.pages[i]
			break
		}
	}

	if d.pageChanged {
		d.pageChanged = false
		page.handler.BeforeRender(d.multiWallet)
	}

	if page.isNavPage {
		d.renderNavPage(page, ctx)
	} else {
		d.renderStandalonePage(page, ctx)
	}
}

func (d *desktop) renderNavPage(page page, ctx *layout.Context) {
	flex := layout.Flex{
		Axis: layout.Horizontal,
	}

	navChild := flex.Rigid(ctx, func() {
		d.renderNavSection(ctx)
	})

	contentChild := flex.Rigid(ctx, func() {
		d.renderContentSection(page, ctx)
	})

	flex.Layout(ctx, navChild, contentChild)
}

func (d *desktop) renderStandalonePage(page page, ctx *layout.Context) {
	page.handler.Render(ctx, d.refreshWindow)
}

func (d *desktop) renderNavSection(ctx *layout.Context) {
	navAreaBounds := image.Point{
		X: navSectionWidth,
		Y: windowHeight * 2,
	}
	helper.PaintArea(ctx, helper.DecredDarkBlueColor, navAreaBounds)

	inset := layout.Inset{
		Top:  unit.Dp(0),
		Left: unit.Dp(0),
	}
	inset.Layout(ctx, func() {
		var stack layout.Stack
		children := make([]layout.StackChild, len(d.pages))

		currentPositionTop := float32(0)
		navButtonHeight := float32(30)

		for index, page := range d.pages {
			children[index] = stack.Rigid(ctx, func() {
				inset := layout.Inset{
					Top:   unit.Dp(currentPositionTop),
					Right: unit.Dp(navSectionWidth),
				}

				c := ctx.Constraints
				inset.Layout(ctx, func() {
					ctx.Constraints.Width.Min = navAreaBounds.X

					for page.button.Clicked(ctx) {
						d.changePage(page.name)
					}
					widgets.LayoutNavButton(page.button, page.label, d.theme, ctx)
					ctx.Constraints = c
				})
			})
			currentPositionTop += navButtonHeight
		}

		stack.Layout(ctx, children...)
	})
}

func (d *desktop) renderContentSection(page page, ctx *layout.Context) {
	inset := layout.Inset{
		Left:  unit.Dp(-113),
		Right: unit.Dp(10),
		Top:   unit.Dp(4),
	}

	inset.Layout(ctx, func() {
		if d.multiWallet.IsSyncing() {
			d.syncer.Render(ctx)
		} else {
			page.handler.Render(ctx, d.refreshWindow)
		}
	})
}

func (d *desktop) refreshWindow() {
	d.window.Invalidate()
}

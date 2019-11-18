package gio

import (
	"fmt"
	"image"
	"log"
	"os"

	gioapp "gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/gio/giolog"
	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
)

type (
	desktop struct {
		window         *gioapp.Window
		displayName    string
		pages          []navPage
		standalonePages map[string]standalonePageHandler
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

	multiWallet, shouldCreateOrRestoreWallet, shouldPromptForPass, err := LoadWallet(appDataDir, netType) 
	if err != nil {
		// todo show error in UI
		giolog.Log.Errorf(err.Error())
		return
	}

	app.multiWallet = multiWallet 
	if shouldCreateOrRestoreWallet {
		app.currentPage = "welcome"
	} else if shouldPromptForPass {
		app.currentPage = "passphrase"
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
	// set standalone page
	d.standalonePages = getStandalonePages(d.multiWallet, d.theme)

	// set navPages
	d.pages = getNavPages()
	if len(d.pages) > 0 && d.currentPage == "" {
		d.changePage(d.pages[0].name)
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
	// first check if current page is standalone and render 
	if page, ok := d.standalonePages[d.currentPage]; ok {
		d.renderStandalonePage(page, ctx)
	} else {
		var page navPage
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

		d.renderNavPage(page, ctx)
	}
}

func (d *desktop) renderNavPage(page navPage, ctx *layout.Context) {
	flex := layout.Flex{
		Axis: layout.Horizontal,
	}

	navChild := flex.Rigid(ctx, func() {
		d.renderNavSection(ctx)
	})

	contentChild := flex.Rigid(ctx, func() {
		inset := layout.Inset{
			Left: unit.Dp(-353),
		}
		inset.Layout(ctx, func(){
			d.renderContentSection(page, ctx)
		})
	})

	flex.Layout(ctx, navChild, contentChild)
}

func (d *desktop) renderStandalonePage(page standalonePageHandler, ctx *layout.Context) {
	gioapp.Size(unit.Dp(200), unit.Dp(windowHeight))
	
	inset := layout.UniformInset(unit.Dp(10))
	inset.Layout(ctx, func(){
		page.Render(ctx, d.refreshWindow, d.changePage)
	})
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
		Right: unit.Dp(float32(navAreaBounds.X)),
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
				}

				c := ctx.Constraints
				inset.Layout(ctx, func() {
					ctx.Constraints.Width.Min = navAreaBounds.X
					
					for page.button.Clicked(ctx) {
						d.changePage(page.name)
					}

					widgets.LayoutNavButton(page.button, page.label, d.theme, ctx)
				})
				ctx.Constraints = c
			})
			currentPositionTop += navButtonHeight
		}
		stack.Layout(ctx, children...)
	})
}

func (d *desktop) renderContentSection(page navPage, ctx *layout.Context) {
	if d.multiWallet.IsSyncing() {
		d.syncer.Render(ctx)
	} else {
		page.handler.Render(ctx, d.refreshWindow)
	}
}

func (d *desktop) refreshWindow() {
	d.window.Invalidate()
}

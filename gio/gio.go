package gio

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"

	gioapp "gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/op/paint"
	"gioui.org/widget/material"

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
		appDisplayName string
		multiWallet    *dcrlibwallet.MultiWallet
		syncer         *Syncer

		logo           material.Image
	}
)

const (
	windowWidth  = 520
	windowHeight = 500

	navSectionWidth = 120

	logoPath = "../../gio/assets/decred.png"
)

func LaunchUserInterface(appDisplayName, appDataDir, netType string) {
	logger, err := dcrlibwallet.RegisterLogger("GIOL")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Launch error - cannot register logger: %v", err)
		return
	}
	giolog.UseLogger(logger)

	// initialize theme 
	helper.Initialize()

	app := &desktop{
		currentPage: "overview",
	}

	theme := helper.GetTheme()
	imageOp, err := app.loadLogo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Launch error - cannot load logo: %v", err)
		return
	}
	app.logo = theme.Image(imageOp)
	app.logo.Scale = 0.95

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

	app.syncer = NewSyncer(app.multiWallet, app.refreshWindow)
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

func (d *desktop) loadLogo() (paint.ImageOp, error) {
	logoByte, err := os.Open(logoPath)
	if err != nil {
		return paint.ImageOp{}, err
	}

	src, _, err := image.Decode(logoByte) 
	if err != nil {
		return paint.ImageOp{}, err
	}

	return paint.NewImageOp(src), nil
}

func (d *desktop) prepareHandlers() {
	// set standalone page
	d.standalonePages = getStandalonePages(d.multiWallet)

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

	d.currentPage = pageName
	d.pageChanged = true
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
	d.refreshWindow()
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
	windowBounds := image.Point{
		X: windowWidth * 2,
		Y: windowHeight * 2,
	}
	helper.PaintArea(ctx, helper.BackgroundColor, windowBounds)

	inset := layout.UniformInset(unit.Dp(25))
	inset.Layout(ctx, func(){
		d.logo.Layout(ctx)
	})

	inset = layout.Inset{
		Top: unit.Dp(65),
		Left: unit.Dp(25),
		Right: unit.Dp(25),
	}
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
					page.button.Draw(ctx, widgets.AlignMiddle, func(){
						d.changePage(page.name)
					})
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

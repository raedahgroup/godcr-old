package gio

import (
	"context"
	"image"
	"log"

	//"gioui.org/ui"
	gioapp "gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"

	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
)

type (
	Desktop struct {
		window           *gioapp.Window
		pages            []page
		currentPage      string
		pageChanged      bool
		theme            *helper.Theme
		syncer           *Syncer
		walletMiddleware app.WalletMiddleware
		settings         *config.Settings
	}
)

const (
	windowWidth  = 450
	windowHeight = 350

	navSectionWidth = 120
)

func LaunchApp(ctx context.Context, walletMiddleware app.WalletMiddleware, settings *config.Settings) {
	theme := helper.NewTheme()

	desktop := &Desktop{
		theme:            theme,
		walletMiddleware: walletMiddleware,
		settings:         settings,
		syncer:           NewSyncer(theme),
		currentPage:      "overview",
	}
	desktop.prepareHandlers()

	go func() {
		desktop.window = gioapp.NewWindow(
			gioapp.Size(unit.Dp(windowWidth), unit.Dp(windowHeight)),
			gioapp.Title(app.DisplayName),
		)

		if err := desktop.renderLoop(); err != nil {
			log.Fatal(err)
		}
	}()

	// start syncing in background
	go desktop.syncer.startSyncing(walletMiddleware, desktop.refreshWindow)

	// run app
	gioapp.Main()
}

func (d *Desktop) prepareHandlers() {
	pages := getPages()
	d.pages = make([]page, len(pages))

	for index, page := range pages {
		d.pages[index] = page

		if index == 0 {
			d.changePage(page.name)
		}
	}
}

func (d *Desktop) changePage(pageName string) {
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

func (d *Desktop) renderLoop() error {
	ctx := &layout.Context{
		Queue: d.window.Queue(),
	}

	for {
		e := <-d.window.Events()
		switch e := e.(type) {
		case gioapp.DestroyEvent:
			return e.Err
		case gioapp.FrameEvent:
			ctx.Reset(&e.Config, e.Size)
			d.render(ctx)
			e.Frame(ctx.Ops)
		}
	}
}

func (d *Desktop) render(ctx *layout.Context) {
	var page page
	for i := range d.pages {
		if d.pages[i].name == d.currentPage {
			page = d.pages[i]
			break
		}
	}

	if d.pageChanged {
		d.pageChanged = false
		page.handler.BeforeRender(d.walletMiddleware, d.settings)
	}

	if page.isNavPage {
		d.renderNavPage(page, ctx)
	} else {
		d.renderStandalonePage(page, ctx)
	}
}

func (d *Desktop) renderNavPage(page page, ctx *layout.Context) {
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

func (d *Desktop) renderStandalonePage(page page, ctx *layout.Context) {
	page.handler.Render(ctx, d.refreshWindow)
}

func (d *Desktop) renderNavSection(ctx *layout.Context) {
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
				ctx.Constraints.Width.Min = 270
				ctx.Constraints.Width.Max = 270

				inset.Layout(ctx, func() {
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

func (d *Desktop) renderContentSection(page page, ctx *layout.Context) {
	inset := layout.Inset{
		Left:  unit.Dp(-113),
		Right: unit.Dp(10),
		Top:   unit.Dp(8),
	}

	inset.Layout(ctx, func() {
		if d.syncer.isDoneSyncing() {
			page.handler.Render(ctx, d.refreshWindow)
		} else {
			d.syncer.Render(ctx)
		}
	})
}

func (d *Desktop) refreshWindow() {
	d.window.Invalidate()
}

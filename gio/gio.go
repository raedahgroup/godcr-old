package gio

import (
	"image"
	"log"

	"gioui.org/gesture"
	"gioui.org/ui"
	gioapp "gioui.org/ui/app"
	"gioui.org/ui/layout"

	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widget"
)

type (
	Desktop struct {
		window      *gioapp.Window
		pages       []page
		navClickers []gesture.Click
		currentPage *page
		pageChanged bool
	}
)

const (
	appName      = "GoDcr"
	windowWidth  = 570
	windowHeight = 400

	navSectionWidth = 150
)

func LaunchApp() {
	desktop := &Desktop{}
	desktop.prepareHandlers()

	helper.Init()
	widget.Init()

	go func() {
		desktop.window = gioapp.NewWindow(
			gioapp.WithWidth(ui.Dp(windowWidth)),
			gioapp.WithHeight(ui.Dp(windowHeight)),
			gioapp.WithTitle(appName),
		)

		if err := desktop.renderLoop(); err != nil {
			log.Fatal(err)
		}
	}()

	gioapp.Main()
}

func (d *Desktop) prepareHandlers() {
	pages := getPages()
	d.pages = make([]page, len(pages))
	d.navClickers = make([]gesture.Click, len(pages))

	for index, page := range pages {
		d.pages[index] = page

		if index == 0 {
			d.changePage(page.name)
		}
	}
}

func (d *Desktop) changePage(pageName string) {
	if d.currentPage != nil && d.currentPage.name == pageName {
		return
	}

	if d.currentPage != nil && d.currentPage.name == pageName {
		return
	}

	for _, page := range d.pages {
		if page.name == pageName {
			d.currentPage = &page
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
		case gioapp.UpdateEvent:
			ctx.Reset(&e.Config, layout.RigidConstraints(e.Size))
			d.render(ctx)
			d.window.Update(ctx.Ops)
		}
	}
}

func (d *Desktop) render(ctx *layout.Context) {
	flex := &layout.Flex{
		Axis: layout.Horizontal,
	}
	flex.Init(ctx)

	navChild := flex.Rigid(func() {
		d.renderNavSection(ctx)
	})

	contentChild := flex.Rigid(func() {
		inset := layout.Inset{
			Left: ui.Dp(navSectionWidth / 2),
			Top: ui.Dp(4),
		}

		inset.Layout(ctx, func() {
			d.renderContentSection(ctx)
		})
	})

	flex.Layout(navChild, contentChild)
}

func (d *Desktop) renderNavSection(ctx *layout.Context) {
	stack := (&layout.Stack{}).Init(ctx)

	navArea := stack.Rigid(func() {
		inset := layout.Inset{
			Left: ui.Dp(0),
			Top:  ui.Dp(0),
		}

		inset.Layout(ctx, func() {
			// paint nav area
			bounds := image.Point{
				X: navSectionWidth,
				Y: windowHeight * 2,
			}
			helper.PaintArea(helper.Theme.Brand, bounds, ctx.Ops)

			positionTop := float32(0)
			for _, page := range d.pages {
				inset := layout.Inset{
					Top: ui.Dp(positionTop),
				}

				inset.Layout(ctx, func() {
					widget.Button(page.label, page.clicker, ctx, func() {
						d.changePage(page.name)
					})
				})
				positionTop += 28
			}
		})
	})

	stack.Layout(navArea)
}

func (d *Desktop) renderContentSection(ctx *layout.Context) {
	if d.pageChanged {
		d.pageChanged = false
		d.currentPage.handler.BeforeRender()
	}

	stack := (&layout.Stack{}).Init(ctx)

	header := stack.Rigid(func() {
		inset := layout.Inset{
			Top:  ui.Dp(0),
			Left: ui.Dp(0),
		}
		inset.Layout(ctx, func() {
			widget.HeadingText(d.currentPage.label, widget.TextAlignLeft, ctx)
		})
	})

	headerLine := stack.Rigid(func() {
		inset := layout.Inset{
			Top: ui.Dp(28),
			Left: ui.Dp(0),
		}

		inset.Layout(ctx, func(){
			bounds := image.Point{
				X: windowWidth - 30,
				Y: 1,
			}
			helper.PaintArea(helper.Theme.Brand, bounds, ctx.Ops)
		})
	})

	stack.Layout(header, headerLine)
}

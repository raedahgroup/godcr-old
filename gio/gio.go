package gio

import (
	"image"
	"log"
	
	//"gioui.org/ui"
	gioapp "gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	//"gioui.org/widget"

	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
	
	
)

type (
	Desktop struct {
		window      *gioapp.Window
		pages       []page
		currentPage *page
		pageChanged bool

		theme *helper.Theme
	}
)

const (
	appName      = "GoDcr"
	windowWidth  = 550
	windowHeight = 450

	navSectionWidth = 120
)

func LaunchApp() {
	desktop := &Desktop{
		theme: helper.NewTheme(),
	}
	desktop.prepareHandlers()
	

	go func() {
		desktop.window = gioapp.NewWindow(
			gioapp.Size(unit.Px(windowWidth), unit.Px(windowHeight)),
			gioapp.Title("GoDcr"),
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
		case gioapp.FrameEvent:
			ctx.Reset(&e.Config, e.Size)
			d.render(ctx)
			e.Frame(ctx.Ops)
		}
	}
}

func (d *Desktop) render(ctx *layout.Context) {
	ctx.Constraints.Width.Min = 0
	ctx.Constraints.Height.Min = 0
	
	flex := layout.Flex{
		Axis: layout.Horizontal,
	}

	navChild := flex.Rigid(ctx, func(){
		d.renderNavSection(ctx)
	})

	contentChild := flex.Rigid(ctx, func(){

	})

	flex.Layout(ctx, navChild, contentChild)
}

func (d *Desktop) renderNavSection(ctx *layout.Context) {
	navAreaBounds := image.Point{
		X: navSectionWidth,
		Y: windowHeight * 2,
	}
	helper.PaintArea(ctx, helper.DecredDarkBlueColor, navAreaBounds)


	inset := layout.Inset{
		Top: unit.Sp(0),
		Left: unit.Sp(0),
		Right: unit.Px(navSectionWidth),
		Bottom: unit.Px(0),
	}
	inset.Layout(ctx, func(){
		var stack layout.Stack 
		children := make([]layout.StackChild, len(d.pages))

		currentPositionTop := float32(0)
		navButtonHeight := float32(30)

		
		for index, page := range d.pages {
			children[index] = stack.Rigid(ctx, func(){
				inset := layout.Inset{
					Top: unit.Sp(currentPositionTop),
					Right: unit.Dp(navSectionWidth),
				}	

				c := ctx.Constraints 
				ctx.Constraints.Width.Min = 270 
				ctx.Constraints.Width.Max = 270

				inset.Layout(ctx, func(){
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



	

	/**var stack layout.Stack 
	children := make([]layout.StackChild, len(d.pages))

	currentPositionTop := float32(0)
	for index, page := range d.pages {
		children[index] = stack.Expand(ctx, func(){
			inset := layout.Inset{
				Top: unit.Sp(currentPositionTop),
			}
	
			inset.Layout(ctx, func(){
				for page.button.Clicked(ctx) {
					d.changePage(page.name)
				}
				d.theme.Button(page.label).Layout(ctx, page.button)
			})
		})
		currentPositionTop += 32
	}

	stack.Layout(ctx, children...)

	//d.theme.Button(d.pages[0].label).Layout(ctx, d.pages[0].button)
	

	/**var stack layout.Stack 
	navArea := stack.Rigid(ctx, func(){
		inset := layout.Inset{
			Top: unit.Sp(0),
			Left: unit.Sp(0),
		}
		inset.Layout(ctx, func(){
			currentPositionTop := float32(0)
			for _, page := range d.pages {
				inset := layout.Inset{
					Top: unit.Sp(currentPositionTop),
					Left: unit.Sp(0),
				}
				inset.Layout(ctx, func(){
					for page.button.Clicked(ctx) {
						d.changePage(page.name)
					}
					d.theme.Button(page.label).Layout(ctx, page.button)
				})
				currentPositionTop += 32
			}
		})
	})
	stack.Layout(ctx, navArea)**/
}

func (d *Desktop) renderContentSection(ctx *layout.Context) {
	if d.pageChanged {
		d.pageChanged = false
		d.currentPage.handler.BeforeRender()
	}

	var stack layout.Stack 

	/**contentAreaBounds := image.Point{
		X: windowWidth * 2,
		Y: windowHeight * 2,
	}

	helper.PaintArea(ctx, helper.BackgroundColor, contentAreaBounds)**/
	
	

	stack.Layout(ctx)

	/**
	stack := (&layout.Stack{})

	header := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top:  unit.Dp(0),
			Left: unit.Dp(0),
		}
		inset.Layout(ctx, func() {
			//widget.HeadingText(d.currentPage.label, widget.TextAlignLeft, ctx)
		})
	})

	headerLine := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(28),
			Left: unit.Dp(0),
		}

		inset.Layout(ctx, func(){
			/**bounds := image.Point{
				X: windowWidth - 30,
				Y: 1,
			}
			helper.PaintArea(helper.Theme.Brand, bounds, ctx.Ops)*
		})
	})

	stack.Layout(ctx, header, headerLine)**/
}

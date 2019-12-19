package widgets 

import (
	"gioui.org/unit"
	"gioui.org/layout"
	"github.com/raedahgroup/godcr/gio/helper"
)

type (
	tabItem struct {
		label *ClickableLabel
		renderFunc func()
	}

	TabContainer struct {
		items 			[]tabItem
		currentTabIndex int
	}
)

func NewTabContainer() *TabContainer {
	return &TabContainer{
		items: []tabItem{},
		currentTabIndex: 0,
	}
}

func (t *TabContainer) AddTab(label string) *TabContainer {
	item := tabItem{
		label     : NewClickableLabel(label), 
		renderFunc: func(){},
	}
	t.items = append(t.items, item)
	return t
}

func (t *TabContainer) Draw(ctx *layout.Context, renderFuncs ...func(*layout.Context)) {
	t.drawNavSection(ctx)

	// draw current tab content
	inset := layout.Inset{
		Top: unit.Dp(25),
	}
	inset.Layout(ctx, func(){
		// todo make sure number of render funcs match number of tabs
		renderFuncs[t.currentTabIndex](ctx)
	})
}

func (t *TabContainer) drawNavSection(ctx *layout.Context) {
	navTabWidth := ctx.Constraints.Width.Max / len(t.items)

	columns := make([]layout.FlexChild, len(t.items))
	for index, tab := range t.items {
		color := helper.BlackColor
		if t.currentTabIndex == index {
			color = helper.DecredLightBlueColor
		}

		columns[index] = layout.Rigid(func(){
			tab.label.SetWidth(navTabWidth).
				SetSize(13).
				SetColor(color).
				SetAlignment(AlignMiddle).
				Draw(ctx, func(){
					t.currentTabIndex = index
				})
		})
	}
	layout.Flex{Axis: layout.Horizontal}.Layout(ctx, columns...)
}

func (t *TabContainer) drawContentSection(ctx *layout.Context, renderFuncs []func(*layout.Context)) {
	inset := layout.Inset{
		Top: unit.Dp(25),
	}
	inset.Layout(ctx, func(){
		
	})
}
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
	stack := layout.Stack{}

	navSection := stack.Rigid(ctx, func(){
		t.drawNavSection(ctx)
	})

	contentSection := stack.Expand(ctx, func(){
		t.drawContentSection(ctx, renderFuncs)
	})

	stack.Layout(ctx, navSection, contentSection)
}

func (t *TabContainer) drawNavSection(ctx *layout.Context) {
	inset := layout.Inset{
		Top:  unit.Sp(0),
		Left: unit.Sp(0),
	}
	inset.Layout(ctx, func() {
		stack := layout.Stack{}
		tabNavWidth := ctx.Constraints.Width.Max / len(t.items)

		navSection := stack.Rigid(ctx, func(){
			flex := layout.Flex{
				Axis: layout.Horizontal,
			}
			children := make([]layout.FlexChild, len(t.items))
			inset := layout.UniformInset(unit.Dp(0))
			for index, tab := range t.items {
				children[index] = flex.Rigid(ctx, func() {
					color := helper.BlackColor
					if t.currentTabIndex == index {
						color = helper.DecredLightBlueColor
					}
					inset.Layout(ctx, func() {
						tab.label.SetWidth(tabNavWidth).
							SetSize(13).
							SetColor(color).
							Draw(ctx, AlignMiddle, func(){
								t.currentTabIndex = index
							})	
					})
				})
			}
			flex.Layout(ctx, children...)
		})
		borderSection := stack.Rigid(ctx, func(){
			inset := layout.Inset{
				Top: unit.Dp(20),
				Left: unit.Dp(0),
				Right: unit.Dp(0),
			}
			inset.Layout(ctx, func(){
				NewLine().SetHeight(2).SetColor(helper.GrayColor).Draw(ctx)
			})
		})

		stack.Layout(ctx, navSection, borderSection)
	})
}

func (t *TabContainer) drawContentSection(ctx *layout.Context, renderFuncs []func(*layout.Context)) {
	inset := layout.Inset{
		Top: unit.Dp(25),
	}
	inset.Layout(ctx, func(){
		// todo make sure number of render funcs match number of tabs
		renderFuncs[t.currentTabIndex](ctx)
	})
}
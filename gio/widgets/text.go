package widgets

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr/gio/helper"
)

func DisplayErrorText(text string, theme *helper.Theme, ctx *layout.Context) {
	label := theme.H5(text)
	label.Color = helper.DangerColor 

	label.Layout(ctx)
}

func AddCenteredLabel(txt string, theme *helper.Theme, ctx *layout.Context) {
	stack := layout.Stack{}
	lbl := stack.Rigid(ctx, func(){
		layout.Align(layout.Center).Layout(ctx, func() {
			widget.Label{Alignment: text.Middle}.Layout(ctx, theme.Shaper, helper.GetNavFont(), txt)
		})
		
	})
	stack.Layout(ctx, lbl)
}

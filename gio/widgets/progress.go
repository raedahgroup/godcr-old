package widgets


import (
	"image"

	"gioui.org/layout"
	"gioui.org/unit"
	//"gioui.org/widget"

	"github.com/raedahgroup/godcr/gio/helper"
)

func NewProgressBar(progress *int, theme *helper.Theme, ctx *layout.Context) {
	stack := layout.Stack{}

	container := stack.Rigid(ctx, func(){
		inset := layout.UniformInset(unit.Dp(15))
		inset.Layout(ctx, func(){
			containerBounds := image.Point{
				X: ctx.Constraints.Width.Max,
				Y: 20,
			}
			helper.PaintArea(ctx, helper.GrayColor, containerBounds)
			indicatorWidth := float64(*progress)/float64(100) * float64(containerBounds.X)
			if indicatorWidth > float64(containerBounds.X) {
				indicatorWidth = float64(containerBounds.X)
			}
			
			indicatorBounds := image.Point{
				X: int(indicatorWidth),
				Y: 20,
			}
			helper.PaintArea(ctx, helper.SuccessColor, indicatorBounds)
		})
	})

	stack.Layout(ctx, container)
}
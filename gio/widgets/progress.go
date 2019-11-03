package widgets

import (
	"image"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/raedahgroup/godcr/gio/helper"
)

const (
	progressBarHeight = 20
)

func NewProgressBar(progress *int, theme *helper.Theme, ctx *layout.Context) {
	stack := layout.Stack{}

	container := stack.Rigid(ctx, func() {
		inset := layout.UniformInset(unit.Dp(0))
		inset.Layout(ctx, func() {
			containerBounds := image.Point{
				X: ctx.Constraints.Width.Max,
				Y: progressBarHeight,
			}
			helper.PaintArea(ctx, helper.GrayColor, containerBounds)
			// calculate width of indicator with respects to progress bar width
			indicatorWidth := float64(*progress) / float64(100) * float64(containerBounds.X)

			if indicatorWidth > float64(containerBounds.X) {
				indicatorWidth = float64(containerBounds.X)
			}

			indicatorBounds := image.Point{
				X: int(indicatorWidth),
				Y: progressBarHeight,
			}
			helper.PaintArea(ctx, helper.SuccessColor, indicatorBounds)
		})
	})

	stack.Layout(ctx, container)
}

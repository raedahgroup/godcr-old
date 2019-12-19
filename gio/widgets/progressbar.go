package widgets

import (
	"image"
	"image/color"

	"gioui.org/layout"
	//"gioui.org/unit"
	"github.com/raedahgroup/godcr/gio/helper"
)

type (
	ProgressBar struct {
		height          int 
		backgroundColor color.RGBA 
		progressColor   color.RGBA
	}
)

const (
	defaultProgressBarHeight = 20
)

func NewProgressBar() *ProgressBar{
	return &ProgressBar{
		height         : defaultProgressBarHeight,
		backgroundColor: helper.GrayColor,
		progressColor  : helper.SuccessColor,
	}
}

func (p *ProgressBar) SetHeight(height int) *ProgressBar {
	p.height = height
	return p
}

func (p *ProgressBar) SetBackgroundColor(col color.RGBA) *ProgressBar {
	p.backgroundColor = col 
	return p
}

func (p *ProgressBar) SetProgressColor(col color.RGBA) *ProgressBar {
	p.progressColor = col
	return p
}

func (p *ProgressBar) Draw(ctx *layout.Context, progress *float64) {
	layout.Stack{}.Layout(ctx, 
		layout.Stacked(func(){
			containerBounds := image.Point{
				X: ctx.Constraints.Width.Max,
				Y: p.height,
			}
			helper.PaintArea(ctx, p.backgroundColor, containerBounds)
			// calculate width of indicator with respects to progress bar width
			indicatorWidth := float64(*progress) / float64(100) * float64(containerBounds.X)
		
			if indicatorWidth > float64(containerBounds.X) {
				indicatorWidth = float64(containerBounds.X)
			}

			indicatorBounds := image.Point{
				X: int(indicatorWidth),
				Y: p.height,
			}
			helper.PaintArea(ctx, p.progressColor, indicatorBounds)
		}),
	)
}

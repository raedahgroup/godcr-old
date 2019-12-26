package helper

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/op/paint"
	"gioui.org/layout"
)

var (
	WhiteColor = color.RGBA{255, 255, 255, 255}
	BlackColor = color.RGBA{0, 0, 0, 255}
	GrayColor  = color.RGBA{200, 200, 200, 255}

	DangerColor  = color.RGBA{215, 58, 73, 255}
	SuccessColor = color.RGBA{227, 98, 9, 255}

	DecredDarkBlueColor  = color.RGBA{9, 20, 64, 255}
	DecredLightBlueColor = color.RGBA{41, 112, 255, 255}

	DecredOrangeColor = color.RGBA{237, 109, 71, 255}
	DecredGreenColor  = color.RGBA{46, 214, 161, 255} //color.RGBA{65, 191, 83, 255}

	BackgroundColor = color.RGBA{243, 245, 246, 255}
)

func PaintArea(ctx *layout.Context, color color.RGBA, x int, y int) {
	bounds := image.Point{
		X: x,
		Y: y,
	}
	
	paint.ColorOp{
		Color: color,
	}.Add(ctx.Ops)
	
	
	paint.PaintOp{
		Rect: f32.Rectangle{
			Max: f32.Point{
				X: float32(bounds.X),
				Y: float32(bounds.Y),
			},
		},
	}.Add(ctx.Ops)
}
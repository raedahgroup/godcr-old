package helper

import (
	"image"
	"image/color"

	"gioui.org/ui/f32"
	"gioui.org/ui/paint"
	"gioui.org/ui"
)

var (
	WhiteColor = color.RGBA{255, 255, 255, 255}
	BlackColor = color.RGBA{0, 0, 0, 255}
	GrayColor  = color.RGBA{200, 200, 200, 255}

	DangerColor  = color.RGBA{215, 58, 73, 255}
	SuccessColor = color.RGBA{227, 98, 9, 255}

	DecredDarkBlueColor  = color.RGBA{9, 20, 64, 255}
	DecredLightBlueColor = color.RGBA{112, 203, 255, 255}

	DecredOrangeColor = color.RGBA{237, 109, 71, 255}
	DecredGreenColor  = color.RGBA{65, 191, 83, 255}
)

func PaintArea(material ui.MacroOp, bounds image.Point, ops *ui.Ops) {
	material.Add(ops)
	paint.PaintOp{
		Rect: f32.Rectangle{
			Max: f32.Point{
				X: float32(bounds.X),
				Y: float32(bounds.Y),
			},
		},
	}.Add(ops)
}

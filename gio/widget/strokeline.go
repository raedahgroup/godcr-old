package widget

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
)

// Strokeline is a line widget
// If tabcontainer, width should be <6 height should be the size of height of the window
// Spacing x should be the max button size and sacing y 0
func Strokeline(gtx *layout.Context, lineColor color.RGBA, width, height, spacingX, spacingY float32) {
	line := f32.Rectangle{Max: f32.Point{X: width, Y: height}}
	op.TransformOp{}.Offset(f32.Point{
		X: spacingX,
		Y: spacingY,
	}).Add(gtx.Ops)
	paint.ColorOp{Color: lineColor}.Add(gtx.Ops)
	paint.PaintOp{Rect: line}.Add(gtx.Ops)
}

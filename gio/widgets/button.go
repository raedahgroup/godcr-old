package widgets

import (
	"gioui.org/text"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/unit"
	"gioui.org/op"
	"github.com/raedahgroup/godcr/gio/helper"

	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/op/paint"
)

func LayoutNavButton(button *widget.Button, txt string, theme *helper.Theme, gtx *layout.Context)  {
	col := helper.WhiteColor
	bgcol := theme.Color.Primary
	if !button.Active() {
		col.A = 0xaa
		bgcol.A = 0xaa
	}
	st := layout.Stack{}
	lbl := st.Rigid(gtx, func() {
		fnt := text.Font{
			Size: theme.TextSize.Scale(14.0 / 16.0),
		}

		layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
			paint.ColorOp{Color: col}.Add(gtx.Ops)
			widget.Label{Alignment: text.Middle}.Layout(gtx, theme.Shaper, fnt, txt)
		})
		pointer.RectAreaOp{Rect: image.Rectangle{Max: gtx.Dimensions.Size}}.Add(gtx.Ops)
		button.Layout(gtx)
	})
	bg := st.Expand(gtx, func() {
		rr := float32(gtx.Px(unit.Dp(4)))
		rrect(gtx.Ops,
			float32(gtx.Constraints.Width.Max),
			float32(gtx.Constraints.Height.Max),
			rr, rr, rr, rr,
		)
		fill(gtx, bgcol)
		for _, c := range button.History() {
			drawInk(gtx, c)
		}
	})
	st.Layout(gtx, bg, lbl)
}


func drawInk(gtx *layout.Context, c widget.Click) {
	d := gtx.Now().Sub(c.Time)
	t := float32(d.Seconds())
	const duration = 0.5
	if t > duration {
		return
	}
	t = t / duration
	var stack op.StackOp
	stack.Push(gtx.Ops)
	size := float32(gtx.Px(unit.Dp(700))) * t
	rr := size * .5
	col := byte(0xaa * (1 - t*t))
	ink := paint.ColorOp{Color: color.RGBA{A: col, R: col, G: col, B: col}}
	ink.Add(gtx.Ops)
	op.TransformOp{}.Offset(c.Position).Offset(f32.Point{
		X: -rr,
		Y: -rr,
	}).Add(gtx.Ops)
	rrect(gtx.Ops, float32(size), float32(size), rr, rr, rr, rr)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(size), Y: float32(size)}}}.Add(gtx.Ops)
	stack.Pop()
	op.InvalidateOp{}.Add(gtx.Ops)
}


func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Max, Y: cs.Height.Max}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d, Baseline: d.Y}
}

// https://pomax.github.io/bezierinfo/#circles_cubic.
func rrect(ops *op.Ops, width, height, se, sw, nw, ne float32) {
	w, h := float32(width), float32(height)
	const c = 0.55228475 // 4*(sqrt(2)-1)/3
	var b paint.Path
	b.Begin(ops)
	b.Move(f32.Point{X: w, Y: h - se})
	b.Cube(f32.Point{X: 0, Y: se * c}, f32.Point{X: -se + se*c, Y: se}, f32.Point{X: -se, Y: se}) // SE
	b.Line(f32.Point{X: sw - w + se, Y: 0})
	b.Cube(f32.Point{X: -sw * c, Y: 0}, f32.Point{X: -sw, Y: -sw + sw*c}, f32.Point{X: -sw, Y: -sw}) // SW
	b.Line(f32.Point{X: 0, Y: nw - h + sw})
	b.Cube(f32.Point{X: 0, Y: -nw * c}, f32.Point{X: nw - nw*c, Y: -nw}, f32.Point{X: nw, Y: -nw}) // NW
	b.Line(f32.Point{X: w - ne - nw, Y: 0})
	b.Cube(f32.Point{X: ne * c, Y: 0}, f32.Point{X: ne, Y: ne - ne*c}, f32.Point{X: ne, Y: ne}) // NE
	b.End().Add(ops)
}
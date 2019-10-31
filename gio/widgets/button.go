package widgets

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr/gio/helper"

	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/op/paint"
)

func NewButton() *widget.Button {
	return new(widget.Button)
}

type IconButton struct {
	Text       string
	Background color.RGBA
	Icon       *Icon
	Size       unit.Value
	Padding    unit.Value
	Font       text.Font
	Color      color.RGBA
	Alignment  layout.Axis

	shaper *text.Shaper
}

type Icon struct {
	imgSize int
	img     image.Image
	src     []byte

	op paint.ImageOp
	material.Icon
}

type Themes struct {
	material.Theme
}

// Todo: depending on dark/white theme this should vary
func NewTheme() *Themes {
	t := &Themes{
		Theme: material.Theme{Shaper: font.Default()},
	}
	t.Color.Primary = color.RGBA{63, 81, 181, 255}
	t.Color.Text = color.RGBA{0, 0, 0, 255}
	t.Color.Hint = color.RGBA{187, 187, 187, 255}
	t.TextSize = unit.Sp(16)
	return t
}

func (t *Themes) Button(icon *Icon, colors color.RGBA, texts string, Alignment layout.Axis, tabContainer bool) IconButton {
	var padding int
	if !tabContainer {
		padding = 20
	}

	return IconButton{
		Background: t.Color.Primary,
		Icon:       icon,
		Size:       unit.Dp(56),
		Padding:    unit.Dp(float32(padding)),
		Text:       texts,
		Color:      t.Color.Text,
		Alignment:  Alignment,

		Font:   text.Font{Size: t.TextSize.Scale(14.0 / 16.0)},
		shaper: t.Shaper,
	}
}

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func (b IconButton) Layout(gtx *layout.Context, button *widget.Button) {
	col := b.Color
	bgcol := b.Background
	if !button.Active() {
		col.A = 0xaa
		bgcol.A = 0xaa
	}

	st := layout.Stack{}
	ico := st.Rigid(gtx, func() {
		iconAndLabel := layout.Flex{Axis: b.Alignment, Alignment: layout.Middle}

		icon := iconAndLabel.Rigid(gtx, func() {
			layout.UniformInset(b.Padding).Layout(gtx, func() {
				size := gtx.Px(b.Size) - 2*gtx.Px(b.Padding)
				if b.Icon != nil {
					ico := b.Icon.image(size)
					ico.Add(gtx.Ops)
					paint.PaintOp{
						Rect: f32.Rectangle{Max: toPointF(ico.Size())}, //toRectF(ico.Bounds()),
					}.Add(gtx.Ops)
				}

				gtx.Dimensions = layout.Dimensions{
					Size: image.Point{X: size, Y: size},
				}
			})
			button.Layout(gtx)
		})

		label := iconAndLabel.Rigid(gtx, func() {
			layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
				paint.ColorOp{Color: col}.Add(gtx.Ops)
				widget.Label{Alignment: text.Middle}.Layout(gtx, b.shaper, b.Font, b.Text)
			})
			button.Layout(gtx)

		})
		iconAndLabel.Layout(gtx, icon, label)
	})

	bg := st.Expand(gtx, func() {
		rr := float32(gtx.Px(unit.Dp(4)))
		rrect(gtx.Ops,
			float32(gtx.Constraints.Width.Min), float32(gtx.Constraints.Height.Min),
			rr, rr, rr, rr,
		)
		fill(gtx, bgcol)
		for _, c := range button.History() {
			drawInk(gtx, c)
		}
	})

	pointer.RectAreaOp{Rect: image.Rectangle{Max: gtx.Dimensions.Size}}.Add(gtx.Ops)
	st.Layout(gtx, bg, ico)
}

func LayoutNavButton(button *widget.Button, txt string, theme *helper.Theme, gtx *layout.Context) {
	col := helper.WhiteColor
	bgcol := helper.DecredDarkBlueColor
	if !button.Active() {
		col = color.RGBA{255, 255, 255, 255}
		bgcol = helper.DecredDarkBlueColor
	}
	st := layout.Stack{Alignment: layout.Center}
	hmin := gtx.Constraints.Width.Min
	//vmin := gtx.Constraints.Height.Min
	lbl := st.Rigid(gtx, func() {
		gtx.Constraints.Width.Min = hmin
		gtx.Constraints.Height.Min = 30
		layout.Align(layout.Center).Layout(gtx, func() {
			layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
				paint.ColorOp{Color: col}.Add(gtx.Ops)
				widget.Label{Alignment: text.Middle}.Layout(gtx, theme.Shaper, helper.GetNavFont(), txt)
			})
		})
		pointer.RectAreaOp{Rect: image.Rectangle{Max: gtx.Dimensions.Size}}.Add(gtx.Ops)
		button.Layout(gtx)
	})
	bg := st.Expand(gtx, func() {
		rr := float32(gtx.Px(unit.Dp(0)))
		rrect(gtx.Ops,
			float32(gtx.Constraints.Width.Min),
			float32(gtx.Constraints.Height.Min),
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

func NewIcon(data []byte) (*Icon, error) {
	_, err := iconvg.DecodeMetadata(data)

	if err != nil {
		return nil, err
	}
	return &Icon{src: data}, nil
}

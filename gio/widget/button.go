package widget

import (
	"image"
	"image/color"
	"image/draw"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/iconvg"
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

func (ic *Icon) image(sz int) paint.ImageOp {
	if sz == ic.imgSize {
		return ic.op
	}
	m, _ := iconvg.DecodeMetadata(ic.src)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	// Use white for icons.
	m.Palette[0] = color.RGBA{A: 0xff, R: 0xff, G: 0xff, B: 0xff}
	iconvg.Decode(&ico, ic.src, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	ic.op = paint.NewImageOp(img)
	ic.imgSize = sz
	return ic.op
}

func toRectF(r image.Rectangle) f32.Rectangle {
	return f32.Rectangle{
		Min: f32.Point{X: float32(r.Min.X), Y: float32(r.Min.Y)},
		Max: f32.Point{X: float32(r.Max.X), Y: float32(r.Max.Y)},
	}
}

func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
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

/**
func Init() {
	clickers = make([]gesture.Click, DummyClicker+1)
}

func Button(text string, clicker int, ctx *layout.Context, handler func()) {
	stack := (&layout.Stack{}).Init(ctx)

	var in layout.Inset

	child := stack.Rigid(func() {
		in = layout.Inset{
			Top:  ui.Dp(10),
			Left: ui.Dp(10),
		}
		in.Layout(ctx, func() {
			ButtonText(text, TextAlignCenter, ctx)
		})
	})

	child1 := stack.Rigid(func() {
		ins := layout.Inset{
			Top:  ui.Dp(in.Top.V - float32(4.5)),
			Left: ui.Dp(in.Left.V - float32(10)),
		}
		ins.Layout(ctx, func() {
			letterSize := 9

			xPoint := letterSize * (len(text) + 1)
			xPoint = 90
			yPoint := 30

			helper.Theme.Brand.Add(ctx.Ops)
			paint.PaintOp{
				Rect: f32.Rectangle{
					Max: f32.Point{
						X: float32(xPoint),
						Y: float32(yPoint),
					},
				},
			}.Add(ctx.Ops)

			click := &clickers[clicker]
			pointer.RectAreaOp{
				Rect: image.Rectangle{
					Max: image.Point{
						X: xPoint,
						Y: yPoint,
					},
				},
			}.Add(ctx.Ops)
			for e, ok := click.Next(ctx); ok; e, ok = click.Next(ctx) {
				if e.Type == gesture.TypeClick {
					handler()
				}
			}
			click.Add(ctx.Ops)
		})
	})
	stack.Layout(child1, child)
}

func ConfirmButton(ctx *layout.Context, text, key string, handler func()) {
	size := 35
	dims := image.Point{X: size, Y: size}

	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(size), Y: float32(size)}}}.Add(ctx.Ops)
	ctx.Dimensions = layout.Dimensions{Size: dims}

	/**licks = append(clicks, gesture.Click{})
	click := &clicks[len(clicks)-1]

	pointer.EllipseAreaOp{
		Rect: image.Rectangle{
			Max: ctx.Dimensions.Size,
		},
	}.Add(ctx.Ops)

	for e, ok := click.Next(ctx); ok; e, ok = click.Next(ctx) {
		if e.Type == gesture.TypeClick {
			fmt.Println("dddd")
		}
	}

	click.Add(ctx.Ops)
}

func toRectF(r image.Rectangle) f32.Rectangle {
	return f32.Rectangle{
		Min: f32.Point{X: float32(r.Min.X), Y: float32(r.Min.Y)},
		Max: f32.Point{X: float32(r.Max.X), Y: float32(r.Max.Y)},
	}
}**/

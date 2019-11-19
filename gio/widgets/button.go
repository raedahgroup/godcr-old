package widgets 

import (
	//"fmt"
	"image"
	"image/color"
	"image/draw"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/io/pointer"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/raedahgroup/godcr/gio/helper"
	"golang.org/x/exp/shiny/iconvg"
)

type (
	Button struct {
		button *widget.Button 
		icon   *Icon
		padding unit.Value 
		size 	unit.Value

		text 	   string

		Color      color.RGBA
		Background color.RGBA
	}

	Icon struct {
		imgSize int 
		img 	image.Image 
		src 	[]byte 

		op paint.ImageOp 
		material.Icon
	}
)

const (
	defaultButtonPadding = 10
)


// NewIcon returns a new Icon from IconVG data.
func NewIcon(data []byte) (*Icon, error) {
	_, err := iconvg.DecodeMetadata(data)
	if err != nil {
		return nil, err
	}
	return &Icon{src: data}, nil
}

func NewButton(txt string, icon *Icon) *Button {
	theme := helper.GetTheme()

	btn := &Button{
		button : new(widget.Button),
		icon   : icon,
		text   :  txt,
		padding: unit.Dp(defaultButtonPadding),
		Background: helper.DecredDarkBlueColor,
		Color: helper.WhiteColor,
	}

	if icon != nil {
		btn.padding = unit.Dp(13)
		btn.size = unit.Dp(46)
	}

	btn.Background = theme.Color.Primary 
	btn.Color      = helper.WhiteColor

	return btn
}

func (b *Button) SetPadding(padding int)  *Button {
	b.padding = unit.Dp(float32(padding))
	return b
}

func (b *Button) SetSize(size int) *Button {
	b.size = unit.Dp(float32(size))
	return b
}

func (b *Button) SetBackgroundColor(color color.RGBA) *Button {
	b.Background = color 
	return b
}

func (b *Button) SetColor(color color.RGBA) *Button {
	b.Color = color 
	return b
}

func (b *Button) Draw(ctx *layout.Context, alignment Alignment, onClick func()) {
	for b.button.Clicked(ctx) {
		onClick()
	}

	theme := helper.GetTheme()
	
	if b.icon != nil {
		b.drawIconButton(ctx, theme, alignment)
		return
	}

	col := b.Color
	bgcol := b.Background
	if !b.button.Active() {
		//col = helper.WhiteColor
		//bgcol = helper.DecredDarkBlueColor
	}
	st := layout.Stack{}
	hmin := ctx.Constraints.Width.Min
	vmin := ctx.Constraints.Height.Min
	lbl := st.Rigid(ctx, func() {
		ctx.Constraints.Width.Min = hmin
		ctx.Constraints.Height.Min = vmin
		layout.Align(layout.Center).Layout(ctx, func() {
			layout.UniformInset(b.padding).Layout(ctx, func() {
				paint.ColorOp{Color: col}.Add(ctx.Ops)
				widget.Label{Alignment: text.Middle}.Layout(ctx, theme.Shaper, theme.Fonts.Bold, b.text)
			})
		})
		pointer.RectAreaOp{Rect: image.Rectangle{Max: ctx.Dimensions.Size}}.Add(ctx.Ops)
		b.button.Layout(ctx)
	})
	bg := st.Expand(ctx, func() {
		rr := float32(ctx.Px(unit.Dp(4)))
		rrect(ctx.Ops,
			float32(ctx.Constraints.Width.Min),
			float32(ctx.Constraints.Height.Min),
			rr, rr, rr, rr,
		)
		fill(ctx, bgcol)
		for _, c := range b.button.History() {
			drawInk(ctx, c)
		}
	})
	st.Layout(ctx, bg, lbl)
}

func (b *Button) drawIconButton(ctx *layout.Context, theme *helper.Theme, alignment Alignment) {
	col   := b.Color 
	bgcol := b.Background 
	if !b.button.Active() {
		//col.A = 0xaa
		//bgcol.A = 0xaa
	}

	stack := layout.Stack{}
	
	lbl := stack.Rigid(ctx, func(){
		alignment := getLayoutAlignment(alignment)
		layout.Align(alignment).Layout(ctx, func(){
			stack := layout.Stack{}
			child := stack.Expand(ctx, func(){
				flex := layout.Flex{}

				iconColumn := flex.Rigid(ctx, func(){
					iconSize := ctx.Px(b.size) - 2*ctx.Px(b.padding)
					topInset :=	(float32(iconSize) - 2*b.padding.V) / 2.0

					inset := layout.Inset{
						Right: b.padding,
						Left: unit.Dp(b.padding.V * 1.2),
						Top: unit.Dp(-topInset+10.0),
					}
					inset.Layout(ctx, func(){
						ico := b.icon.image(iconSize)
						ico.Add(ctx.Ops)
						paint.PaintOp{
							Rect: f32.Rectangle{Max: toPointF(ico.Size())}, //toRectF(ico.Bounds()),
						}.Add(ctx.Ops)
					})
				})

				textColumn := flex.Rigid(ctx, func(){
					inset := layout.Inset{
						Left: unit.Dp(b.padding.V + 7.0),
						Top: b.padding,
						Bottom: b.padding,
					}
					inset.Layout(ctx, func(){
						paint.ColorOp{Color: col}.Add(ctx.Ops)
						widget.Label{Alignment: text.Middle}.Layout(ctx, theme.Shaper, theme.Fonts.Regular, b.text)
					})
				})
				flex.Layout(ctx, iconColumn, textColumn)
			})
			stack.Layout(ctx, child)
		})
	})

	clickDimensionSize := ctx.Dimensions.Size 
	

	bg := stack.Expand(ctx, func(){
		ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
		rr := float32(ctx.Px(unit.Dp(4)))
		rrect(ctx.Ops,
			float32(ctx.Constraints.Width.Min),
			float32(ctx.Constraints.Height.Min),
			rr, rr, rr, rr,
		)
		fill(ctx, bgcol)

		clickDimensionSize.X = ctx.Dimensions.Size.X
		pointer.RectAreaOp{Rect: image.Rectangle{Max: clickDimensionSize}}.Add(ctx.Ops)
		b.button.Layout(ctx)

		for _, c := range b.button.History() {
			drawInk(ctx, c)
		}
	})
	stack.Layout(ctx, bg, lbl)
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


func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func drawInk(ctx *layout.Context, c widget.Click) {
	d := ctx.Now().Sub(c.Time)
	t := float32(d.Seconds())
	const duration = 0.5
	if t > duration {
		return
	}
	t = t / duration
	var stack op.StackOp
	stack.Push(ctx.Ops)
	size := float32(ctx.Px(unit.Dp(700))) * t
	rr := size * .5
	col := byte(0xaa * (1 - t*t))
	ink := paint.ColorOp{Color: color.RGBA{A: col, R: col, G: col, B: col}}
	ink.Add(ctx.Ops)
	op.TransformOp{}.Offset(c.Position).Offset(f32.Point{
		X: -rr,
		Y: -rr,
	}).Add(ctx.Ops)
	rrect(ctx.Ops, float32(size), float32(size), rr, rr, rr, rr)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(size), Y: float32(size)}}}.Add(ctx.Ops)
	stack.Pop()
	op.InvalidateOp{}.Add(ctx.Ops)
}

func fill(ctx *layout.Context, col color.RGBA) {
	cs := ctx.Constraints
	d := image.Point{X: cs.Width.Max, Y: cs.Height.Max}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(ctx.Ops)
	paint.PaintOp{Rect: dr}.Add(ctx.Ops)
	ctx.Dimensions = layout.Dimensions{Size: d, Baseline: d.Y}
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

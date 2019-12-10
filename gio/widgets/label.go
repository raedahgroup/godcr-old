package widgets

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/f32"
	"gioui.org/widget/material"
	"github.com/raedahgroup/godcr/gio/helper"
	"golang.org/x/image/math/fixed"
)

type (
	Label struct {
		material.Label
		size int
	}

	Alignment int
)

const (
	AlignLeft Alignment = iota
	AlignMiddle
	AlignRight
)

const (
	NormalLabelHeight = 30
)

func NewLabel(txt string, size ...int) *Label {
	labelSize := 3 
	if len(size) > 0 {
		labelSize = size[0]
	}

	return &Label{
		Label: getLabelWithSize(txt, labelSize),
		size: labelSize,
	}
}

func NewErrorLabel(txt string) *Label {
	l := &Label{
		Label: getLabelWithSize(txt, 4),
		size: 4,
	}

	return l.SetColor(helper.DangerColor)
}

func (l *Label) SetText(txt string) *Label {
	l.Label = getLabelWithSize(txt, l.size)
	return l
}

func (l *Label) SetSize(size int) *Label {
	l.Label = getLabelWithSize(l.Text, size)
	l.size = size
	return l
}

func (l *Label) SetWeight(weight text.Weight) *Label {
	l.Font.Weight = weight
	return l
}

func (l *Label) SetStyle(style text.Style) *Label {
	l.Font.Style = style 
	return l
}

func (l *Label) SetColor(color color.RGBA) *Label {
	l.Label.Color = color
	return l
}

func (l *Label) Draw(ctx *layout.Context, alignment Alignment) {
	l.Label.Alignment = getTextAlignment(alignment)
	l.Label.Layout(ctx)
}


type ClickableLabel struct {
	label *Label 
	clicker helper.Clicker
	width  int
}

func NewClickableLabel(txt string, size ...int) *ClickableLabel {
	labelSize := 2 
	if len(size) > 0 {
		labelSize = size[0]
	}
	
	return &ClickableLabel{
		label  : NewLabel(txt, labelSize).SetColor(helper.DecredDarkBlueColor),
		clicker: helper.NewClicker(),
	}
}


func (c *ClickableLabel) SetText(txt string) *ClickableLabel {
	c.label.SetText(txt)
	return c
}

func (c *ClickableLabel) SetSize(size int) *ClickableLabel {
	c.label.SetSize(size)
	return c
}

func (c *ClickableLabel) SetStyle(style text.Style) *ClickableLabel {
	c.label.SetStyle(style)
	return c
}

func (c *ClickableLabel) SetWeight(weight text.Weight) *ClickableLabel {
	c.label.SetWeight(weight)
	return c
}

func (c *ClickableLabel) SetColor(color color.RGBA) *ClickableLabel {
	c.label.SetColor(color)
	return c
}

func (c *ClickableLabel) SetWidth(width int) *ClickableLabel {
	c.width = width 
	return c
}

func (c *ClickableLabel) Draw(ctx *layout.Context, alignment Alignment, onClick func()) {
	for c.clicker.Clicked(ctx) {
		onClick()
	}

	stack := layout.Stack{}
	child := stack.Expand(ctx, func(){
		if c.width == 0 {
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
		} else {
			ctx.Constraints.Width.Min =  c.width
		}
		

		c.label.Draw(ctx, alignment)
		pointer.RectAreaOp{Rect: image.Rectangle{Max: ctx.Dimensions.Size}}.Add(ctx.Ops)
		c.clicker.Register(ctx)
	})
	stack.Layout(ctx, child)
}


func getTextAlignment(alignment Alignment) text.Alignment {
	switch alignment {
	case AlignMiddle:
		return text.Middle
	case AlignRight:
		return text.End
	default:
		return text.Start
	}
}

func getLayoutAlignment(alignment Alignment) layout.Alignment {
	switch alignment {
	case AlignMiddle:
		return layout.Middle
	case AlignRight:
		return layout.End
	default:
		return layout.Start
	}
}

func getLabelWithSize(txt string, size int) material.Label {
	theme := helper.GetTheme()

	switch size {
	case 1:
		return theme.Caption(txt)
	case 2:
		return theme.Body2(txt)
	case 3:
		return theme.Body1(txt)
	case 4:
		return theme.H6(txt)
	case 5:
		return theme.H5(txt)
	case 6:
		return theme.H4(txt)
	case 7:
		return theme.H3(txt)
	case 8:
		return theme.H2(txt)
	case 9:
		return theme.H1(txt)
	default:
		return theme.Body1(txt)
	}
}	


func TextPadding(lines []text.Line) (padding image.Rectangle) {
	if len(lines) == 0 {
		return
	}
	first := lines[0]
	if d := first.Ascent + first.Bounds.Min.Y; d < 0 {
		padding.Min.Y = d.Ceil()
	}
	last := lines[len(lines)-1]
	if d := last.Bounds.Max.Y - last.Descent; d > 0 {
		padding.Max.Y = d.Ceil()
	}
	if d := first.Bounds.Min.X; d < 0 {
		padding.Min.X = d.Ceil()
	}
	if d := first.Bounds.Max.X - first.Width; d > 0 {
		padding.Max.X = d.Ceil()
	}
	return
}

func LinesDimens(lines []text.Line) layout.Dimensions {
	var width fixed.Int26_6
	var h int
	var baseline int
	if len(lines) > 0 {
		baseline = lines[0].Ascent.Ceil()
		var prevDesc fixed.Int26_6
		for _, l := range lines {
			h += (prevDesc + l.Ascent).Ceil()
			prevDesc = l.Descent
			if l.Width > width {
				width = l.Width
			}
		}
		h += lines[len(lines)-1].Descent.Ceil()
	}
	w := width.Ceil()
	return layout.Dimensions{
		Size: image.Point{
			X: w,
			Y: h,
		},
		Baseline: h - baseline,
	}
}


func Align(align text.Alignment, width fixed.Int26_6, maxWidth int) fixed.Int26_6 {
	mw := fixed.I(maxWidth)
	switch align {
	case text.Middle:
		return fixed.I(((mw - width) / 2).Floor())
	case text.End:
		return fixed.I((mw - width).Floor())
	case text.Start:
		return 0
	default:
		panic(fmt.Errorf("unknown alignment %v", align))
	}
}


func ToRectF(r image.Rectangle) f32.Rectangle {
	return f32.Rectangle{
		Min: f32.Point{X: float32(r.Min.X), Y: float32(r.Min.Y)},
		Max: f32.Point{X: float32(r.Max.X), Y: float32(r.Max.Y)},
	}
}
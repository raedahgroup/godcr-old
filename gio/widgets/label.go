package widgets

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/raedahgroup/godcr/gio/helper"
)

type (
	Label struct {
		material.Label
		size int
	}
)

const (
	AlignLeft = iota
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
	l.Font.Size = unit.Dp(float32(size))
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

func (l *Label) Draw(ctx *layout.Context, alignment int) {
	l.Label.Alignment = getTextAlignment(alignment)
	l.Label.Layout(ctx)
}


type ClickableLabel struct {
	label *Label 
	clicker helper.Clicker
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

func (c *ClickableLabel) Draw(ctx *layout.Context, alignment int, onClick func()) {
	for c.clicker.Clicked(ctx) {
		onClick()
	}

	stack := layout.Stack{}
	child := stack.Rigid(ctx, func(){
		ctx.Constraints.Width.Min = ctx.Constraints.Width.Max

		c.label.Draw(ctx, alignment)
		pointer.RectAreaOp{Rect: image.Rectangle{Max: ctx.Dimensions.Size}}.Add(ctx.Ops)
		c.clicker.Register(ctx)
	})
	stack.Layout(ctx, child)
}


func getTextAlignment(alignment int) text.Alignment {
	switch alignment {
	case AlignMiddle:
		return text.Middle
	case AlignRight:
		return text.End
	default:
		return text.Start
	}
}

func getLayoutAlignment(alignment int) layout.Alignment {
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
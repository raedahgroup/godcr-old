package widgets

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/text"

	"github.com/raedahgroup/godcr/gio/helper"
)

type (
	ClickableLabel struct {
		text      string
		clicker   helper.Clicker
		theme     *helper.Theme
		alignment int
	}
)

const (
	NormalLabelHeight = 15

	AlignLeft = iota
	AlignMiddle
	AlignRight
)

func DisplayErrorText(text string, theme *helper.Theme, ctx *layout.Context) {
	label := theme.H5(text)
	label.Color = helper.DangerColor

	label.Layout(ctx)
}

func CenteredLabel(txt string, theme *helper.Theme, ctx *layout.Context) {
	ctx.Constraints.Width.Min = ctx.Constraints.Width.Max

	label := theme.H6(txt)
	label.Alignment = text.Middle
	label.Layout(ctx)
}

func BoldCenteredLabel(txt string, theme *helper.Theme, ctx *layout.Context) {
	ctx.Constraints.Width.Min = ctx.Constraints.Width.Max

	label := theme.H4(txt)
	label.Alignment = text.Middle
	label.Layout(ctx)
}

func NewClickableLabel(txt string, alignment int, theme *helper.Theme) *ClickableLabel {
	return &ClickableLabel{
		text:      txt,
		alignment: alignment,
		theme:     theme,
		clicker:   helper.NewClicker(),
	}
}

func (c *ClickableLabel) SetText(txt string) {
	c.text = txt
}

func (c *ClickableLabel) Display(onClick func(), ctx *layout.Context) {
	for c.clicker.Clicked(ctx) {
		onClick()
	}

	stack := layout.Stack{}
	child := stack.Rigid(ctx, func() {
		ctx.Constraints.Width.Min = ctx.Constraints.Width.Max

		var alignment text.Alignment
		switch c.alignment {
		case AlignMiddle:
			alignment = text.Middle
		case AlignRight:
			alignment = text.End
		default:
			alignment = text.Start
		}

		label := c.theme.H6(c.text)
		label.Alignment = alignment
		label.Layout(ctx)

		pointer.RectAreaOp{Rect: image.Rectangle{Max: ctx.Dimensions.Size}}.Add(ctx.Ops)
		c.clicker.Register(ctx)
	})
	stack.Layout(ctx, child)
}

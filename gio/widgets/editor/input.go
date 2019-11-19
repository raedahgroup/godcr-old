package editor 

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
)

type (
	Input struct {
		*Editor
		hint   string
	}
)

func NewInput(hint string) *Input { 
	return &Input{
		Editor: new(Editor),
		hint  : hint,
	}
}

func (i *Input) SetMask(char string) *Input {
	i.setMask(char)
	return i
}

func (i *Input) Draw(ctx *layout.Context) {
	theme := helper.GetTheme()

	var stack op.StackOp 
	stack.Push(ctx.Ops)
	var macro op.MacroOp 
	macro.Record(ctx.Ops)
	paint.ColorOp{
		Color: helper.GrayColor,
	}.Add(ctx.Ops)
	widgets.NewLabel(i.hint, 3).SetColor(helper.GrayColor).Draw(ctx, widgets.AlignLeft)
	macro.Stop()
	if w := ctx.Dimensions.Size.X; ctx.Constraints.Width.Min < w {
		ctx.Constraints.Width.Min = w
	}
	if h := ctx.Dimensions.Size.Y; ctx.Constraints.Height.Min < h {
		ctx.Constraints.Height.Min = h
	}
	i.Layout(ctx, theme.Shaper, theme.Fonts.Regular)
	if i.Len() > 0 {
		paint.ColorOp{
			Color: helper.BlackColor,
		}.Add(ctx.Ops)
		i.PaintText(ctx)
	} else {
		macro.Add(ctx.Ops)
	}
	paint.ColorOp{
		Color: helper.BlackColor,
	}.Add(ctx.Ops)
	i.PaintCaret(ctx)
	stack.Pop()
}
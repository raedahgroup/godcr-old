package widget

import (


	"gioui.org/widget"
)


func NewButton() *widget.Button {
	return new(widget.Button)
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

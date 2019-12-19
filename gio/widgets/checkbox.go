package widgets 

import (
	//"image"

	//"gioui.org/f32"
	//"gioui.org/io/pointer"
	"gioui.org/layout"
	//"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	//"github.com/raedahgroup/godcr/gio/helper"
)



type (
	Checkbox struct {
		isChecked bool
		icon   	  *Icon
		padding   unit.Value
		size      unit.Value
		button    *widget.Button 
	}
)

func NewCheckbox() *Checkbox {
	return &Checkbox{
		isChecked: false,
		icon: 	   NavigationCheckIcon,
		button:    new(widget.Button),
		padding:   unit.Dp(5),
		size:      unit.Dp(26),
	}
}

func (c *Checkbox) IsChecked() bool {
	return c.isChecked
}

func (c *Checkbox) toggleCheckState() {
	if c.isChecked {
		c.isChecked = false 
		return
	}
	c.isChecked = true
}

func (c *Checkbox) Draw(ctx *layout.Context) {
	/**for c.button.Clicked(ctx) {
		c.toggleCheckState()
	}


	st := layout.Stack{}
	hmin := ctx.Constraints.Width.Min
	vmin := ctx.Constraints.Height.Min
	ico := st.Rigid(ctx, func() {
		ctx.Constraints.Width.Min = hmin
		ctx.Constraints.Height.Min = vmin
		layout.Align(layout.Center).Layout(ctx, func() {
			layout.UniformInset(c.padding).Layout(ctx, func() {
				size := ctx.Px(c.size) - 2*ctx.Px(c.padding)
				if c.isChecked {
					ico := c.icon.image(size)
					ico.Add(ctx.Ops)
					paint.PaintOp{
						Rect: f32.Rectangle{Max: toPointF(ico.Size())}, //toRectF(ico.Bounds()),
					}.Add(ctx.Ops)
				}
				ctx.Dimensions = layout.Dimensions{
					Size: image.Point{X: size, Y: size},
				}
			})
		})
		pointer.EllipseAreaOp{Rect: image.Rectangle{Max: ctx.Dimensions.Size}}.Add(ctx.Ops)
		c.button.Layout(ctx)
	})
	bgcol := helper.DecredGreenColor
	if !c.isChecked {
		bgcol = helper.WhiteColor
	}
	bg := st.Expand(ctx, func() {
		ctx.Constraints.Width.Min = 36

		size   := float32(ctx.Constraints.Width.Min) 
		rr     := float32(size) * .5
		Rrect(ctx.Ops,
			float32(ctx.Constraints.Width.Min),
			float32(ctx.Constraints.Width.Min),
			rr, rr, rr, rr,
		)
		Fill(ctx, helper.DecredGreenColor)

		layout.Align(layout.Center).Layout(ctx, func(){
			layout.UniformInset(unit.Dp(1)).Layout(ctx, func(){
				ctx.Constraints.Width.Min = 34

				mainSize   := float32(ctx.Constraints.Width.Min) 
				mainRadius := float32(mainSize) * .5
		
				Rrect(ctx.Ops,
					mainSize,
					mainSize,
					mainRadius, mainRadius, mainRadius, mainRadius,
				)
				Fill(ctx, bgcol)
				for _, c := range c.button.History() {
					drawInk(ctx, c)
				}
			})
		})
	})
	st.Layout(ctx, bg, ico)**/
}
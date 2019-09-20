package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type ClickableBox struct {
	*widget.Box

	OnTapped func() `json:"-"`
}

// Tapped is called when users click on the icon
func (c *ClickableBox) Tapped(_ *fyne.PointEvent) {
	c.OnTapped()
}

// TappedSecondary is called when users right click on the icon
func (c *ClickableBox) TappedSecondary(_ *fyne.PointEvent) {
	// handle secondary tapped (right click)
}

func (c *ClickableBox) CreateRenderer() fyne.WidgetRenderer {
	return widget.Renderer(c.Box)
}

func (c *ClickableBox) Refresh() {
	object := fyne.CurrentApp().Driver().CanvasForObject(c)
	object.Refresh(c)
}

func NewClickableBox(box *widget.Box, OnTapped func()) *ClickableBox {
	icon := box

	clickable := &ClickableBox{icon, OnTapped}
	return clickable
}

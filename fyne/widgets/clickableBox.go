package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type ClickableBox struct {
	*widget.Box

	disable  bool
	OnTapped func() `json:"-"`
}

func NewClickableBox(box *widget.Box, OnTapped func()) *ClickableBox {
	icon := box

	clickable := &ClickableBox{icon, false, OnTapped}
	return clickable
}

func (c *ClickableBox) Disable() {
	c.disable = true
}

func (c *ClickableBox) Enable() {
	c.disable = false
}

func (c *ClickableBox) Disabled() bool {
	return c.disable
}

// Tapped is called when users click on the icon
func (c *ClickableBox) Tapped(_ *fyne.PointEvent) {
	if c.disable == true {
		return
	}

	if c.OnTapped == nil {
		return
	}

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

	if object == nil {
		return
	}

	object.Refresh(c)
}

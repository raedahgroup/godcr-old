package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type ClickableWidget struct {
	fyne.Widget

	disable  bool
	OnTapped func() `json:"-"`
}

func NewClickableWidget(object fyne.Widget, OnTapped func()) *ClickableWidget {
	clickable := &ClickableWidget{object, false, OnTapped}
	return clickable
}

func (c *ClickableWidget) Disable() {
	c.disable = true
}

func (c *ClickableWidget) Enable() {
	c.disable = false
}

func (c *ClickableWidget) Disabled() bool {
	return c.disable
}

// Tapped is called when users click on the icon
func (c *ClickableWidget) Tapped(_ *fyne.PointEvent) {
	if c.disable == true {
		return
	}

	if c.OnTapped == nil {
		return
	}

	c.OnTapped()
}

// TappedSecondary is called when users right click on the icon
func (c *ClickableWidget) TappedSecondary(_ *fyne.PointEvent) {
	// handle secondary tapped (right click)
}

func (c *ClickableWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.Renderer(c.Widget)
}

func (c *ClickableWidget) Refresh() {
	c.Widget.Refresh()
}

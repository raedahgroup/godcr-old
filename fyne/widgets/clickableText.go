package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

type ClickableText struct {
	*canvas.Text

	OnTapped func() `json:"-"`
}

func NewClickableText(box *canvas.Text, OnTapped func()) *ClickableText {
	icon := box
	clickable := &ClickableText{icon, OnTapped}
	return clickable
}

// Tapped is called when users click on the icon
func (c *ClickableText) Tapped(_ *fyne.PointEvent) {
	if c.OnTapped == nil {
		return
	}

	c.OnTapped()
}

// TappedSecondary is called when users right click on the icon
func (c *ClickableText) TappedSecondary(_ *fyne.PointEvent) {
	// handle secondary tapped (right click)
}

func (c *ClickableText) CreateRenderer() fyne.WidgetRenderer {
	return c.CreateRenderer()
}

func (c *ClickableText) Refresh() {
	object := fyne.CurrentApp().Driver().CanvasForObject(c)
	object.Refresh(c)
}

package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type ClickableIcon struct {
	*widget.Icon
	OnTapped func() `json:"-"`
}

func (c *ClickableIcon) Tapped(_ *fyne.PointEvent) {
	c.OnTapped()
}

func (c *ClickableIcon) TappedSecondary(_ *fyne.PointEvent) {
	// handle secondary tapped (right click)
}

func (c *ClickableIcon) CreateRenderer() fyne.WidgetRenderer {
	return widget.Renderer(c.Icon)
}

func (c *ClickableIcon) SetIcon(res fyne.Resource) {
	c.Icon.SetResource(res)
	c.Refresh()
}

func (c *ClickableIcon) Refresh() {
	object := fyne.CurrentApp().Driver().CanvasForObject(c)
	object.Refresh(c)
}

func NewClickableIcon(res fyne.Resource, OnTapped func()) *ClickableIcon {
	icon := widget.NewIcon(res)
	clickable := &ClickableIcon{icon, OnTapped}
	return clickable
}

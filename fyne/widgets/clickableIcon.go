package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/widget"
)

type ClickableIcon struct {
	*widget.Icon

	shadow   *widget.Icon
	OnTapped func() `json:"-"`
}

// Tapped is called when users click on the icon
func (c *ClickableIcon) Tapped(_ *fyne.PointEvent) {
	c.OnTapped()
}

// TappedSecondary is called when users right click on the icon
func (c *ClickableIcon) TappedSecondary(_ *fyne.PointEvent) {
	// handle secondary tapped (right click)
}

// MouseIn is called when a desktop pointer enters the widget
// when mouse is hovering clickable icon, shadowed image should be shown
func (c *ClickableIcon) MouseIn(*desktop.MouseEvent) {
	if c.shadow == nil {
		return
	}

	resource := c.Resource
	c.Icon.SetResource(c.shadow.Resource)
	c.shadow.Resource = resource
	c.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
// When mouse isn't hovering clickable icon, shadowed image wont be shown
func (c *ClickableIcon) MouseOut() {
	if c.shadow == nil {
		return
	}

	resource := c.Resource
	c.Icon.SetResource(c.shadow.Resource)
	c.shadow.Resource = resource
	c.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (c *ClickableIcon) MouseMoved(*desktop.MouseEvent) {
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

func NewClickableIcon(res fyne.Resource, shadow fyne.Resource, OnTapped func()) *ClickableIcon {
	icon := widget.NewIcon(res)
	var shadowIcon *widget.Icon
	if shadow != nil {
		shadowIcon = widget.NewIcon(res)
	}
	clickable := &ClickableIcon{icon, shadowIcon, OnTapped}
	return clickable
}

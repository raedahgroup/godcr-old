package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type Container struct {
	*fyne.Container
}

func NewHBoxContainer() *Container {
	return &Container{
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
		),
	}
}

func NewFixedGridLayout(width, height int) *Container {
	size := fyne.NewSize(width, height)

	return &Container{
		fyne.NewContainerWithLayout(
			layout.NewFixedGridLayout(size),
		),
	}
}

func NewScrollableContainer(width, height int, content *Box) *widget.ScrollContainer {
	container := widget.NewScrollContainer(content.Box)
	container.Resize(fyne.NewSize(200, 200))

	return container
}

func (c *Container) AddChildContainer(container *Container) *Container {
	c.Container.AddObject(container.Container)
	return c
}

func (c *Container) AddBox(b *Box) {
	c.AddObject(b.Box)
}

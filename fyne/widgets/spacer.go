package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
)

func NewVSpacer(height int) fyne.CanvasObject {
	space := fyne.NewSize(0, height)
	return fyne.NewContainerWithLayout(layout.NewFixedGridLayout(space), layout.NewSpacer())
}

func NewHSpacer(width int) fyne.CanvasObject {
	space := fyne.NewSize(width, 0)
	return fyne.NewContainerWithLayout(layout.NewFixedGridLayout(space), layout.NewSpacer())
}

package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func CenterObject(object fyne.CanvasObject, bordered bool) fyne.CanvasObject {
	if bordered {
		return NewVBox(layout.NewSpacer(), object, layout.NewSpacer())
	}

	return widget.NewVBox(layout.NewSpacer(), object, layout.NewSpacer())
}

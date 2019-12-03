package widgets

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

// ImplementBorder gets minimum size of all objects.
// Note: this doesnt consider objects that have been resized.
func ImplementBorder(color color.Color, padding fyne.Size, objects ...fyne.CanvasObject) *canvas.Rectangle {
	border := canvas.NewRectangle(color) //theme.BackgroundColor())
	border.StrokeColor = color

	var minSize fyne.Size

	for _, object := range objects {
		minSize = minSize.Add(object.MinSize())
	}

	minSize = minSize.Add(padding)

	border.SetMinSize(minSize)
	return border
}

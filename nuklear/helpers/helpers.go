package helpers

import (
	"image"

	"github.com/aarzilli/nucular"
)

func DrawLink(window *nucular.Window, text string, handler func(text string, window *nucular.Window)) {
	style := window.Master().Style()
	style.Selectable.Hover.Data.Color = contentBackgroundColor
	style.Selectable.HoverActive.Data.Color = contentBackgroundColor
	style.Selectable.Normal.Data.Color = contentBackgroundColor
	style.Selectable.Pressed.Data.Color = contentBackgroundColor
	style.Selectable.Normal.Data.Color = contentBackgroundColor
	style.Selectable.TextPressed = navBackgroundColor
	style.Selectable.TextNormal = navBackgroundColor
	style.Selectable.TextHover = navBackgroundColor

	style.Selectable.Padding = image.Point{0, 0}

	window.Master().SetStyle(style)

	val := false
	if window.SelectableLabel(text, "LC", &val) {
		handler(text, window)
	}
}

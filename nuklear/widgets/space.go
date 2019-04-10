package widgets

import (
	"github.com/aarzilli/nucular/rect"
)

// AddSpacing uses `SpaceBegin` and `LayoutSpacePushScaled`
// to create a rectangle of the specified `width` and `height` on the window.
// Use 0 for `width` or `height` to use all available width or height.
func (window *Window) AddSpacing(width, height int) {
	lineArea := window.Row(height).SpaceBegin(0)

	if width == 0 {
		width = lineArea.W
	}
	if height == 0 {
		height = lineArea.H
	}

	window.LayoutSpacePushScaled(rect.Rect{
		X: 0,
		Y: 0,
		W: width,
		H: height,
	})
}

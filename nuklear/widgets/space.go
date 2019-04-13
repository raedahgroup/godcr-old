package widgets

import "github.com/aarzilli/nucular/rect"

// AddHorizontalSpace uses `SpaceBegin` and `LayoutSpacePushScaled`
// to create a rectangle of the max available width and specified `height` on the window.
func (window *Window) AddHorizontalSpace(height int) {
	lineArea := window.Row(height).SpaceBegin(0)
	window.LayoutSpacePushScaled(rect.Rect{
		X: 0,
		Y: 0,
		W: lineArea.W,
		H: height,
	})
}

package widgets

import (
	"image/color"

	"github.com/aarzilli/nucular/rect"
)

// AddHorizontalLine uses `SpaceBegin` and `LayoutSpacePushScaled`
// to create a rectangle of max width and the specified `height` on the window.
// The rectangle background color is then set to the specified `color`.
// After the rectangle is drawn using a `NoScrollGroupWindow`,
// the `window` background color is reset to what it was before the line was drawn.
func (window *Window) AddHorizontalLine(height int, color color.RGBA) {
	lineArea := window.Row(height).SpaceBegin(0)
	window.LayoutSpacePushScaled(rect.Rect{
		X: 0,
		Y: 0,
		W: lineArea.W,
		H: height,
	})

	windowStyle := window.Master().Style()
	currentWindowBackground := windowStyle.GroupWindow.FixedBackground.Data.Color
	windowStyle.GroupWindow.FixedBackground.Data.Color = color

	NoScrollGroupWindow("line", window.Window, func(w *Window) {
		// reset background
		windowStyle.GroupWindow.FixedBackground.Data.Color = currentWindowBackground
	})
}

// AddVerticalLine uses `SpaceBegin` and `LayoutSpacePushScaled`
// to create a rectangle of the specified `width` and max height on the window.
// The rectangle background color is then set to the specified `color`.
// After the rectangle is drawn using a `NoScrollGroupWindow`,
// the `window` background color is reset to what it was before the line was drawn.
func (window *Window) AddVerticalLine(width int, color color.RGBA) {
	lineArea := window.Row(0).SpaceBegin(0)
	window.LayoutSpacePushScaled(rect.Rect{
		X: 0,
		Y: 0,
		W: width,
		H: lineArea.H,
	})

	windowStyle := window.Master().Style()
	currentWindowBackground := windowStyle.GroupWindow.FixedBackground.Data.Color
	windowStyle.GroupWindow.FixedBackground.Data.Color = color

	NoScrollGroupWindow("line", window.Window, func(w *Window) {
		// reset background
		windowStyle.GroupWindow.FixedBackground.Data.Color = currentWindowBackground
	})
}

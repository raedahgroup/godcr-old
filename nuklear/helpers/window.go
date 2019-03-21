package helpers

import (
	"image"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
)

type Window struct {
	*nucular.Window
}

var contentArea rect.Rect

func NewWindow(title string, w *nucular.Window, flags nucular.WindowFlags) *Window {
	if nw := w.GroupBegin(title, flags); nw != nil {
		return &Window{
			nw,
		}
	}
	return nil
}

func (w *Window) DrawHeader(title string) {
	w.Row(40).Dynamic(1)

	// style header group
	style := w.Master().Style()
	style.GroupWindow.FixedBackground.Data.Color = colorWhite
	style.Font = PageHeaderFont
	w.Master().SetStyle(style)

	if group := w.GroupBegin(title, nucular.WindowNoScrollbar); group != nil {
		group.Row(40).Dynamic(1)

		// set padding
		style = group.Master().Style()
		style.GroupWindow.Padding = image.Point{18, 15}
		group.Master().SetStyle(style)

		if paddedWindow := group.GroupBegin("padded window", nucular.WindowNoScrollbar); paddedWindow != nil {
			paddedWindow.Row(24).Dynamic(1)
			paddedWindow.Label(title, "LC")
			paddedWindow.GroupEnd()
		}

		// reset padding
		style = group.Master().Style()
		style.GroupWindow.Padding = noPadding
		group.Master().SetStyle(style)

		group.GroupEnd()
	}

	// reset style
	style.GroupWindow.FixedBackground.Data.Color = colorContentBackground
	style.GroupWindow.Padding = noPadding
	style.Font = PageContentFont
	w.Master().SetStyle(style)
}

func (w *Window) ContentWindow(title string) *Window {
	w.Row(0).Dynamic(1)
	w.Style()
	return NewWindow(title, w.Window, 0)
}

func (w *Window) SetErrorMessage(message string) {
	w.Row(300).Dynamic(1)
	w.LabelWrap(message)
}

func (w *Window) Style() {
	style := w.Master().Style()
	style.GroupWindow.Padding = image.Point{20, 20}
	style.GroupWindow.Header.Normal.Data.Color = colorAccent
	style.Font = PageContentFont

	w.Master().SetStyle(style)
}

func (w *Window) End() {
	w.GroupEnd()
}

func SetContentArea(area rect.Rect) {
	contentArea = area
}

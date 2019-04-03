package helpers

import (
	"image"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
	nstyle "github.com/aarzilli/nucular/style"
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
	w.Row(50).Dynamic(1)
	bounds := rect.Rect{
		X: contentArea.X,
		Y: contentArea.Y,
		W: contentArea.W,
		H: 80,
	}

	_, out := w.Custom(nstyle.WidgetStateActive)
	if out != nil {
		out.FillRect(bounds, 0, whiteColor)
	}

	bounds.Y += 25
	bounds.X += 30

	out.DrawText(bounds, title, PageHeaderFont, colorTable.ColorText)
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
	style.GroupWindow.Header.Normal.Data.Color = secondaryColor
	style.Font = PageContentFont

	w.Master().SetStyle(style)
}

func (w *Window) End() {
	w.GroupEnd()
}

func SetContentArea(area rect.Rect) {
	contentArea = area
}

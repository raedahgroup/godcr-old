package widgets

import (
	"image"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/nuklear/styles"
)

type Window struct {
	*nucular.Window
}

var pageContentPadding = image.Point{10, 10}

func GroupWindow(uniqueWindowTitle string, parentWindow *nucular.Window, flags nucular.WindowFlags, windowReady func(*Window)) {
	if nw := parentWindow.GroupBegin(uniqueWindowTitle, flags); nw != nil {
		window := &Window{nw}
		defer window.DoneAddingWidgets()
		windowReady(window)
	}
}

func NoScrollGroupWindow(uniqueWindowTitle string, parentWindow *nucular.Window, windowReady func(*Window)) {
	GroupWindow(uniqueWindowTitle, parentWindow, nucular.WindowNoScrollbar, windowReady)
}

func DefaultGroupWindow(uniqueWindowTitle string, parentWindow *nucular.Window, windowReady func(*Window)) {
	GroupWindow(uniqueWindowTitle, parentWindow, 0, windowReady)
}

func PageContentWindow(pageTitle string, parentWindow *nucular.Window, windowReady func(contentWindow *Window)) {
	NoScrollGroupWindow(pageTitle + "-page", parentWindow, func(pageWindow *Window) {
		pageWindow.SetPageTitle(pageTitle)
		pageWindow.ContentWindow(pageTitle + "-page-content", windowReady)
	})
}

func (window *Window) ContentWindow(uniqueWindowTitle string, windowReady func(*Window)) {
	// window should take available height and width
	window.Row(0).Dynamic(1)

	// add padding and set font
	style := window.Master().Style()
	style.GroupWindow.Padding = pageContentPadding
	style.Font = styles.PageContentFont
	window.Master().SetStyle(style)

	// create group window
	DefaultGroupWindow(uniqueWindowTitle, window.Window, windowReady)
}

func (window *Window) SetErrorMessage(message string) {
	window.Row(300).Dynamic(1)
	window.LabelWrap(message)
}

func (window *Window) SetPageTitle(title string) {
	// change window font and draw page title label
	masterWindow := window.Master()
	currentStyle := masterWindow.Style()
	currentFont := currentStyle.Font
	currentStyle.Font = styles.PageHeaderFont
	masterWindow.SetStyle(currentStyle)

	// draw page title using label
	window.Row(30).Dynamic(1)
	window.Label(title, LeftCenterAlign)

	// reset font
	currentStyle.Font = currentFont
	masterWindow.SetStyle(currentStyle)
}

func (window *Window) DoneAddingWidgets() {
	window.GroupEnd()
}

package widgets

import (
	"image"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"golang.org/x/image/font"
	"github.com/aarzilli/nucular/rect"
)

type Window struct {
	*nucular.Window
}

func GroupWindow(uniqueWindowTitle string, parentWindow *nucular.Window, flags nucular.WindowFlags, windowReady func(*Window)) {
	if nw := parentWindow.GroupBegin(uniqueWindowTitle, flags); nw != nil {
		window := &Window{nw}
		windowReady(window)
		window.DoneAddingWidgets()
	}
}

func NoScrollGroupWindow(uniqueWindowTitle string, parentWindow *nucular.Window, windowReady func(*Window)) {
	GroupWindow(uniqueWindowTitle, parentWindow, nucular.WindowNoScrollbar, windowReady)
}

func ScrollableGroupWindow(uniqueWindowTitle string, parentWindow *nucular.Window, windowReady func(*Window)) {
	GroupWindow(uniqueWindowTitle, parentWindow, 0, windowReady)
}

func PageContentWindowWithTitle(pageTitle string, parentWindow *nucular.Window, windowReady func(contentWindow *Window)) {
	PageContentWindowWithTitleAndPadding(pageTitle, parentWindow, 0, 0, windowReady)
}

func PageContentWindowWithTitleAndPadding(pageTitle string, parentWindow *nucular.Window, xPadding, yPadding int, windowReady func(contentWindow *Window)) {
	NoScrollGroupWindow(pageTitle+"-page", parentWindow, func(pageWindow *Window) {
		pageWindow.Master().Style().GroupWindow.Padding = image.Point{X: xPadding, Y: yPadding}
		pageWindow.SetPageTitle(pageTitle)
		pageWindow.PageContentWindow(pageTitle+"-page-content", windowReady)
	})
}

func (window *Window) PageContentWindow(uniqueWindowTitle string, windowReady func(*Window)) {
	// create a rect for this page content window to prevent styles from spilling into other windows
	pageContentArea := window.Row(0).SpaceBegin(1)
	window.LayoutSpacePushScaled(rect.Rect{
		X: 0,
		Y: 0,
		W: pageContentArea.W,
		H: pageContentArea.H,
	})

	window.Master().Style().Font = styles.PageContentFont

	// create group window
	ScrollableGroupWindow(uniqueWindowTitle, window.Window, windowReady)
}

func (window *Window) SetFont(font font.Face) {
	window.Master().Style().Font = font
}

func (window *Window) UseFontAndResetToPrevious(font fontFace, fontReadyForUse func()) {
	currentFont := window.Master().Style().Font
	if currentFont != font {
		window.SetFont(font)
		defer window.SetFont(currentFont)
	}

	fontReadyForUse()
}

func (window *Window) SetPageTitle(title string) {
	window.AddLabelWithFont(title, LeftCenterAlign, styles.PageHeaderFont)
}

func (window *Window) DisplayErrorMessage(errorMessage string) {
	window.AddWrappedLabelWithColor(errorMessage, styles.DecredOrangeColor)
}

func (window *Window) DoneAddingWidgets() {
	window.GroupEnd()
}

package widgets

import (
	"image"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"golang.org/x/image/font"
)

const defaultPageContentPadding = 10

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

func PageContentWindowDefaultPadding(pageTitle string, parentWindow *nucular.Window, windowReady func(contentWindow *Window)) {
	PageContentWindowWithPadding(pageTitle, parentWindow, defaultPageContentPadding, 0, windowReady)
}

func PageContentWindowWithPadding(pageTitle string, parentWindow *nucular.Window, xPadding, yPadding int, windowReady func(contentWindow *Window)) {
	NoScrollGroupWindow(pageTitle+"-page", parentWindow, func(pageWindow *Window) {
		pageWindow.Master().Style().GroupWindow.Padding = image.Point{X: xPadding, Y: yPadding}

		pageWindow.AddSpacing(0, defaultPageContentPadding)
		pageWindow.SetPageTitle(pageTitle)
		pageWindow.AddSpacing(0, defaultPageContentPadding)

		pageWindow.PageContentWindow(pageTitle+"-page-content", xPadding, yPadding, windowReady)
	})
}

func (window *Window) PageContentWindow(uniqueWindowTitle string, xPadding, yPadding int, windowReady func(*Window)) {
	// create a rect for this page content window to prevent styles from spilling into other windows
	pageContentArea := window.Row(0).SpaceBegin(1)
	window.LayoutSpacePushScaled(rect.Rect{
		X: 0,
		Y: 0,
		W: pageContentArea.W,
		H: pageContentArea.H,
	})

	window.Master().Style().Font = styles.PageContentFont
	window.Master().Style().GroupWindow.Padding = image.Point{X: xPadding, Y: yPadding}

	// create group window
	ScrollableGroupWindow(uniqueWindowTitle, window.Window, windowReady)
}

func (window *Window) Font() font.Face {
	return window.Master().Style().Font
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
	window.AddWrappedLabelWithColor("Error: "+errorMessage, styles.DecredOrangeColor)
}

func (window *Window) DoneAddingWidgets() {
	window.GroupEnd()
}

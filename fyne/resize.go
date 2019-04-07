/*
This file contains workarounds to achieve proper sizing of window content on the fyne app.
This should be replaced with a standard window resize event listener from the fyne library.
*/

package fyne

import (
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type resizeListenerData = struct {
	lastWindowSize    fyne.Size
	everySecondTicker *time.Ticker
}

func (app *fyneApp) listenForWindowResizeEvents() {
	resizeListener := &resizeListenerData{
		lastWindowSize:    app.mainWindow.Content().Size(),
		everySecondTicker: time.NewTicker(100 * time.Millisecond),
	}

	go func() {
		for {
			select {
			case <-resizeListener.everySecondTicker.C:
				newSize := app.mainWindow.Content().Size()
				previousSize := resizeListener.lastWindowSize
				if newSize.Height != previousSize.Height || newSize.Width != previousSize.Width {
					resizeListener.lastWindowSize = newSize
					app.resizeScrollableContainer()
				}
			}
		}
	}()
}

// resizeScrollableContainer ensures that
// - the content of each page is wrapped in scrollable container
// - the scrollable container takes the maximum space available
//
// The idea is, if the content size is bigger than the maximum space available,
// scroll bars become visible and more of the content can be seen by scrolling.
func (app *fyneApp) resizeScrollableContainer() {
	// calculate the maximum available width and height to use for scroll container
	windowSize := app.mainWindow.Content().Size()
	scrollAreaWidth := windowSize.Width - menuSectionWidth - menuSectionPageSectionSeparation
	scrollAreaHeight := windowSize.Height - app.pageTitle.Size().Height
	scrollAreaSize := fyne.NewSize(scrollAreaWidth, scrollAreaHeight)

	// use calculated max size to layout the scrollable container
	scrollContainerLayout := layout.NewFixedGridLayout(scrollAreaSize)
	scrollableContainer := fyne.NewContainerWithLayout(scrollContainerLayout, widget.NewScrollContainer(app.pageContent))

	// must clear items and re-add otherwise the added content will not display
	app.pageSectionOnRight.Children = []fyne.CanvasObject{}
	app.pageSectionOnRight.Append(app.pageTitle)
	app.pageSectionOnRight.Append(scrollableContainer)
}

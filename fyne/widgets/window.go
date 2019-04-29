package widgets

import (
	"fyne.io/fyne"
)

type Window struct {
	fyne.App
	main      fyne.Window
	container *Container
}

func NewWindow(title string, app fyne.App) *Window {
	return &Window{
		App:  app,
		main: app.NewWindow(title),
	}
}

func (window *Window) Content() fyne.CanvasObject {
	return window.main.Content()
}

func (window *Window) RefreshContainer(container *Container) {
	window.main.Canvas().Refresh(container.Container)
}

func (window *Window) Render(container *Container) {
	window.main.SetContent(container.Container)
	window.main.ShowAndRun()
}

func (window *Window) Close() {
	window.App.Quit()
}

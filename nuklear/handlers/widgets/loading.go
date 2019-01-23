package widgets

import (
	"github.com/aarzilli/nucular"
)

type LoadingWidget struct {
	window *nucular.Window
}

func NewLoadingWidget() Widget {
	return &LoadingWidget{}
}

func (l *LoadingWidget) BeforeRender(window *nucular.Window) {
	l.window = window
}

func (l *LoadingWidget) Render(finishHandler func()) {
	l.window.Row(30).Dynamic(1)
	l.window.Label("Fetching data...", "LC")
}

func (l *LoadingWidget) AfterRender() {
}

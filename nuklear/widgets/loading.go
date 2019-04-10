package widgets

import (
	"github.com/aarzilli/nucular"
)

type LoadingWidget struct {
}

func NewLoadingWidget() Widget {
	return &LoadingWidget{}
}

func (l *LoadingWidget) Render(window *nucular.Window) {
	window.Row(30).Dynamic(1)
	window.Label("Fetching data...", "LC")
}

package widgets

import (
	"github.com/aarzilli/nucular"
)

type Widget interface {
	Render(window *nucular.Window)
}

func ShowLoadingWidget(window *nucular.Window) {
	NewLoadingWidget().Render(window)
}

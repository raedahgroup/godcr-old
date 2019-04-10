package widgets

import "github.com/aarzilli/nucular/label"

const buttonHeight = 20

func (window *Window) AddButton(buttonText string, buttonClickFunc func()) {
	buttonLabel := label.TA(buttonText, CenterAlign)

	window.Row(buttonHeight).Dynamic(1)
	if window.Button(buttonLabel, false) {
		buttonClickFunc()
	}
}


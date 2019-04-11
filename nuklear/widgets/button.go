package widgets

import (
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
)

const buttonPadding = 10
const buttonHeight = 30
const bigButtonHeight = 40

// AddButton adds a button to the window that uses just the width required to display the button text and a standard height
func (window *Window) AddButton(buttonText string, buttonClickFunc func()) {
	buttonLabel := label.TA(buttonText, CenterAlign)
	textWidth := nucular.FontWidth(window.Master().Style().Font, buttonText)

	window.Row(buttonHeight).Static(textWidth + buttonPadding)
	if window.Button(buttonLabel, false) {
		buttonClickFunc()
	}
}

// AddBigButton adds a button to the window that stretches to fill the available width and uses a greater height value
func (window *Window) AddBigButton(buttonText string, buttonClickFunc func()) {
	buttonLabel := label.TA(buttonText, CenterAlign)

	window.Row(bigButtonHeight).Dynamic(1)
	if window.Button(buttonLabel, false) {
		buttonClickFunc()
	}
}

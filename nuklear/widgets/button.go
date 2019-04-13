package widgets

import (
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
)

const ButtonHeight = 30
const bigButtonHeight = 40

// AddButton adds a button to the window that uses just the width required to display the button text and a standard height
func (window *Window) AddButton(buttonText string, buttonClickFunc func()) {
	buttonWidth := window.ButtonWidth(buttonText)
	window.Row(ButtonHeight).Static(buttonWidth)

	buttonLabel := label.TA(buttonText, CenterAlign)
	if window.Button(buttonLabel, false) {
		buttonClickFunc()
	}
}

// AddButton adds a button to the window that uses just the width required to display the button text and a standard height
func (window *Window) AddButtonToCurrentRow(buttonText string, buttonClickFunc func()) {
	buttonLabel := label.TA(buttonText, CenterAlign)
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

func (window *Window) ButtonWidth(buttonText string) int {
	textWidth := nucular.FontWidth(window.Master().Style().Font, buttonText)
	buttonPadding := window.Master().Style().Button.Padding.X * 2
	return textWidth + buttonPadding
}

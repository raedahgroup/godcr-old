package widgets

import (
	"math"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	f "golang.org/x/image/font"
	"image/color"
)

const (
	LeftCenterAlign = "LC"
	CenterAlign     = "CC"
)

type fontFace f.Face

// AddLabel adds a single line label to the window. The label added does not wrap.
func (window *Window) AddLabel(text string, align label.Align) {
	window.AddLabelWithFont(text, align, window.Master().Style().Font)
}

// AddLabelWithFont adds a single line label to the window. The label added does not wrap.
func (window *Window) AddLabelWithFont(text string, align label.Align, font fontFace) {
	singleLineHeight := nucular.FontHeight(font) + 1
	window.DrawLabel(text, singleLineHeight, align, font)
}

func (window *Window) DrawLabel(text string, height int, align label.Align, font fontFace) {
	if height < 20 {
		height = 20 // seems labels will not be drawn if row height is less than 20
	}

	window.UseFontAndResetToPrevious(font, func() {
		window.Row(height).Dynamic(1)
		window.Label(text, align)
	})
}

// AddWrappedLabel adds a label to the window.
// The label added wraps it's text and assumes the height required to display all it's text.
func (window *Window) AddWrappedLabel(text string) {
	window.AddWrappedLabelWithFont(text, window.Master().Style().Font)
}

// AddWrappedLabel adds a label to the window.
// The label added wraps it's text and assumes the height required to display all it's text.
func (window *Window) AddWrappedLabelWithColor(text string, color color.RGBA) {
	font := window.Master().Style().Font
	textHeight := window.WrappedLabelTextHeight(text, font)
	window.DrawWrappedLabel(text, textHeight, font)
}

func (window *Window) AddWrappedLabelWithFont(text string, font fontFace) {
	textHeight := window.WrappedLabelTextHeight(text, font)
	window.DrawWrappedLabel(text, textHeight, font)
}

func (window *Window) WrappedLabelTextHeight(text string, font fontFace) int {
	textWidth := nucular.FontWidth(font, text)

	nLines := math.Ceil(float64(textWidth) / float64(window.LayoutAvailableWidth()))
	singleLineHeight := nucular.FontHeight(font) + 1

	return int(nLines) * singleLineHeight + 20 // seems labels will not be drawn if row height is not way higher than necessary
}

func (window *Window) DrawWrappedLabel(text string, height int, font fontFace) {
	window.UseFontAndResetToPrevious(font, func() {
		window.Row(height + 20).Dynamic(1)
		window.LabelWrap(text)
	})
}

func (window *Window) DrawWrappedLabelWithColor(text string, height int, color color.RGBA, font fontFace) {
	window.UseFontAndResetToPrevious(font, func() {
		window.Row(height).Dynamic(1)
		window.LabelWrapColored(text, color)
	})
}

package widgets

import (
	"math"

	"github.com/aarzilli/nucular"
)

const (
	LeftCenterAlign = "LC"
	CenterAlign = "CC"
)

func (window *Window) AddLabel(text string) {
	// measure size required to display text using the current font
	font := window.Master().Style().Font
	textWidth := nucular.FontWidth(font, text)

	nLines := math.Ceil(float64(textWidth) / float64(window.Bounds.W))
	singleLineHeight := nucular.FontHeight(font) + 1
	textHeight := int(nLines) * singleLineHeight

	window.Row(textHeight).Dynamic(1)
	window.LabelWrap(text)
}

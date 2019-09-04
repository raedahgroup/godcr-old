package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func NewItalicizedLabel(text string) *widget.Label {
	return widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
}

package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func sendPageContent() fyne.CanvasObject {
	return widget.NewLabelWithStyle("Send", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
}

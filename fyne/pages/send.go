package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func SendPageContent() fyne.CanvasObject {
	return widget.NewLabelWithStyle("Send", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
}

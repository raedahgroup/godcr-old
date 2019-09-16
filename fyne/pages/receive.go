package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func receivePageContent() fyne.CanvasObject {
	return widget.NewLabelWithStyle("Receive", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
}

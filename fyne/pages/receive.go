package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func ReceivePageContent() fyne.CanvasObject {
	return widget.NewLabelWithStyle("Receive", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
}

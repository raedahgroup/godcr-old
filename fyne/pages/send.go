package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type sendPageDynamicData struct {
	fromAccountSelect *widget.Select
	toAccountSelect   *widget.Select
}

func sendPageContent() fyne.CanvasObject {
	return widget.NewLabelWithStyle("Send", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
}

package history

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func PageContent() fyne.CanvasObject {
	return widget.NewLabelWithStyle("History", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
}

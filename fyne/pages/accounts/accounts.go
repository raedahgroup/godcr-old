package accounts

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func PageContent() fyne.CanvasObject {
	return widget.NewLabelWithStyle("Accounts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
}

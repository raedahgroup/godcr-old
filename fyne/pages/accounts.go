package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func AccountsPageContent() fyne.CanvasObject {
	return widget.NewLabelWithStyle("Accounts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
}

package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// todo: display overview page (include sync progress UI elements)
// todo: register sync progress listener on overview page to update sync progress views
func OverviewPageContent() fyne.CanvasObject {
	return widget.NewLabelWithStyle("Overview", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
}

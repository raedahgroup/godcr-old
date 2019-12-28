package staking

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func PageContent() fyne.CanvasObject {
	return widget.NewLabelWithStyle("Staking", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
}

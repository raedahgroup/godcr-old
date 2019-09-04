package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func StakingPageContent() fyne.CanvasObject {
	return widget.NewLabelWithStyle("Staking", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
}

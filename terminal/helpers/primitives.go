package helpers

import "github.com/rivo/tview"

func CenterAlignedTextView(text string) *tview.TextView {
	return NewTextView(text, tview.AlignCenter)
}

func NewTextView(text string, alignment int) *tview.TextView {
	return tview.NewTextView().SetTextAlign(alignment).SetText(text)
}

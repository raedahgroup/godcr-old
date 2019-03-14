package primitives

import "github.com/rivo/tview"

func WordWrappedTextView(text string) *tview.TextView {
	return NewCenterAlignedTextView(text).
		SetWordWrap(true).
		SetWrap(true)
}

func NewCenterAlignedTextView(text string) *tview.TextView {
	return NewTextView(text, tview.AlignCenter)
}

func NewTextView(text string, alignment int) *tview.TextView {
	return tview.NewTextView().SetTextAlign(alignment).SetText(text)
}

package primitives

import "github.com/rivo/tview"

func WordWrappedTextView(text string) *tview.TextView {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetWordWrap(true).
		SetWrap(true).
		SetText(text)
}

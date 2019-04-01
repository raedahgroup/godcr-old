package primitives

import (
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/rivo/tview"
)

func TitleTextView(text string) *tview.TextView {
	return NewLeftAlignedTextView(text).
		SetTextColor(helpers.PageHeaderColor)
}

func WordWrappedTextView(text string) *tview.TextView {
	return NewLeftAlignedTextView(text).
		SetWordWrap(true).
		SetWrap(true)
}

func NewLeftAlignedTextView(text string) *tview.TextView {
	return NewTextView(text, tview.AlignLeft)
}

func NewCenterAlignedTextView(text string) *tview.TextView {
	return NewTextView(text, tview.AlignCenter)
}

func NewTextView(text string, alignment int) *tview.TextView {
	return tview.NewTextView().SetTextAlign(alignment).SetText(text)
}

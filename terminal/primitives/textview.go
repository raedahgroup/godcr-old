package primitives

import "github.com/rivo/tview"

type TextView struct {
	*tview.TextView
	text string
	border bool
}

// SetText sets the text of this text view to the provided string. Previously
// contained text will be removed.
// Wrapper around `tview.TextView.SetText`
func (t *TextView) SetText(text string) *TextView {
	t.text = text
	t.TextView.SetText(text)
	return t
}

// GetText returns the text displayed in this textview
func (t *TextView) GetText() string {
	return t.text
}

// SetBorder sets the flag indicating whether or not the box should have a border.
// Wrapper around `tview.Box.SetBorder`.
func (t *TextView) SetBorder(show bool) *TextView {
	t.border = show
	t.TextView.SetBorder(show)
	return t
}

// HasBorder returns true if this textview is set to have borders
func (t *TextView) HasBorder() bool {
	return t.border
}

func WordWrappedTextView(text string) *TextView {
	t := NewCenterAlignedTextView(text)
	t.TextView.SetWordWrap(true).SetWrap(true)
	return t
}

func NewCenterAlignedTextView(text string) *TextView {
	return NewTextView(text, tview.AlignCenter)
}

func NewTextView(text string, alignment int) *TextView {
	return &TextView{
		TextView: tview.NewTextView().SetTextAlign(alignment).SetText(text),
		text:text,
	}
}

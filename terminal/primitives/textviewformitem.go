package primitives

import (
	"math"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type TextViewFormItem struct {
	*TextView
	fieldWidth    int
	fieldHeight   int
	autosize      bool
	autosizeWidth int
}

func NewTextViewFormItem(textView *TextView, fieldWidth, fieldHeight int, autosize bool, autosizeWidth int) *TextViewFormItem {
	return &TextViewFormItem{
		TextView:      textView,
		fieldWidth:    fieldWidth,
		fieldHeight:   fieldHeight,
		autosize:      autosize,
		autosizeWidth: autosizeWidth,
	}
}

// GetFieldHeight satisfies `primitives.FormItem` interface
func (item *TextViewFormItem) GetFieldHeight() int {
	if item.autosize {
		availableWidth := item.autosizeWidth
		if item.autosizeWidth <= 0 {
			_, _, availableWidth, _ = item.GetTextView().GetInnerRect()
		}

		textWidth := tview.StringWidth(item.GetText())
		fieldHeight := math.Ceil(float64(textWidth) / float64(availableWidth))

		if item.HasBorder() {
			fieldHeight += 2
		}

		return int(fieldHeight)
	}

	return item.fieldHeight
}

// GetLabel satisfies `tview.FormItem` interface
func (item *TextViewFormItem) GetLabel() string {
	return ""
}

// SetFormAttributes satisfies `tview.FormItem` interface
func (item *TextViewFormItem) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	return item
}

// GetFieldWidth satisfies `tview.FormItem` interface
func (item *TextViewFormItem) GetFieldWidth() int {
	return item.fieldWidth
}

// SetFinishedFunc satisfies `tview.FormItem` interface
func (item *TextViewFormItem) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	// todo confirm this step
	item.GetTextView().SetDoneFunc(handler)
	return item
}

func (item *TextViewFormItem) GetTextView() *TextView {
	return item.TextView
}

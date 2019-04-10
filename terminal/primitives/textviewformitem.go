package primitives

import (
	"math"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type TextViewFormItem struct {
	*TextView
	label       string
	labelWidth  int
	labelColor  tcell.Color
	fieldWidth  int
	fieldHeight int
	autosize    bool
}

func NewTextViewFormItem(textView *TextView, fieldWidth, fieldHeight int, autosize bool) *TextViewFormItem {
	item := &TextViewFormItem{
		TextView:    textView,
		fieldWidth:  fieldWidth,
		fieldHeight: fieldHeight,
		autosize:    autosize,
	}

	return item
}

func (item *TextViewFormItem) CalculateFieldSize(maxWidth int) {
	if item.autosize {
		borderWidth, borderHeight := 0, 0
		if item.HasBorder() {
			borderWidth, borderHeight = 2, 2
		}

		textWidth := tview.StringWidth(item.GetText())
		fieldHeight := math.Ceil(float64(textWidth) / float64(maxWidth-borderWidth))

		item.fieldHeight = int(fieldHeight) + borderHeight
		item.fieldWidth = maxWidth
	}
}

// GetFieldHeight satisfies `primitives.FormItem` interface
func (item *TextViewFormItem) GetFieldHeight() int {
	return item.fieldHeight
}

func (item *TextViewFormItem) SetLabel(label string) *TextViewFormItem {
	item.label = label
	return item
}

// GetLabel satisfies `tview.FormItem` interface
func (item *TextViewFormItem) GetLabel() string {
	return item.label
}

// SetFormAttributes satisfies `tview.FormItem` interface
func (item *TextViewFormItem) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	item.labelWidth = labelWidth
	item.labelColor = labelColor
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

func (item *TextViewFormItem) Draw(screen tcell.Screen) {
	// call textview.Draw() directly if there's no label to display
	if item.label == "" {
		item.TextView.Draw(screen)
		return
	}

	// Prepare
	x, y, width, height := item.GetInnerRect()
	if height < 1 || width < 1 {
		return
	}
	rightLimit := x + width

	// Draw label.
	if item.labelWidth > 0 {
		labelWidth := item.labelWidth
		if labelWidth > width {
			labelWidth = width
		}
		tview.Print(screen, item.label, x, y, labelWidth, tview.AlignLeft, item.labelColor)
		x += labelWidth
	} else {
		_, drawnWidth := tview.Print(screen, item.label, x, y, rightLimit-x, tview.AlignLeft, item.labelColor)
		x += drawnWidth
	}

	// Draw embedded textview using adjusted x pos and width.
	fieldWidth := item.fieldWidth
	if fieldWidth == 0 {
		fieldWidth = math.MaxInt32
	}
	if fieldWidth > rightLimit-x {
		fieldWidth = rightLimit - x
	}

	item.TextView.SetRect(x, y, fieldWidth, height)
	item.TextView.Draw(screen)
}

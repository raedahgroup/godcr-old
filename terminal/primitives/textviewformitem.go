package primitives

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type TextViewFormItem struct {
	tview.Primitive
}

func NewTextViewFormItem(textView *tview.TextView) *TextViewFormItem {
	textView.SetRect(0, 0, 0, 400)
	return &TextViewFormItem{
		Primitive: textView,
	}
}

func (item *TextViewFormItem) GetLabel() string {
	return "hey"
}

func (item *TextViewFormItem) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	return item
}

func (item *TextViewFormItem) GetFieldWidth() int {
	return 0
}

func (item *TextViewFormItem) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	// todo confirm this step
	item.GetTextView().SetDoneFunc(handler)
	return item
}

func (item *TextViewFormItem) GetTextView() *tview.TextView {
	if item.Primitive == nil {
		return nil
	}
	return item.Primitive.(*tview.TextView)
}

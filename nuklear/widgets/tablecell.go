package widgets

import (
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"image/color"
)

type TableCell interface {
	// Render will be called when this table cell item is ready to be added to the window.
	// The table will have constructed the row and column needed to hold this item.
	Render(window *Window)

	// MinWidth returns the min width required to draw this item.
	MinWidth(window *Window) int
}

type LabelTableCell struct {
	text     string
	align    label.Align
	color    color.RGBA
	colorSet bool
}

func NewLabelTableCell(text string, align label.Align) *LabelTableCell {
	return &LabelTableCell{
		text:  text,
		align: align,
	}
}

func NewColoredLabelTableCell(text string, align label.Align, color color.RGBA) *LabelTableCell {
	return &LabelTableCell{
		text:     text,
		align:    align,
		color:    color,
		colorSet: true,
	}
}

func (label *LabelTableCell) Render(window *Window) {
	if label.colorSet {
		window.LabelColored(label.text, label.align, label.color)
	} else {
		window.Label(label.text, label.align)
	}
}

func (label *LabelTableCell) MinWidth(window *Window) int {
	return nucular.FontWidth(window.Font(), label.text)
}

type EditTableCell struct {
	nucular.TextEditor
	width int
}

func NewEditTableCell(editor nucular.TextEditor, width int) *EditTableCell {
	return &EditTableCell{
		editor,
		width,
	}
}

func (edit *EditTableCell) Render(window *Window) {
	edit.Edit(window.Window)
}

func (edit *EditTableCell) MinWidth(_ *Window) int {
	return edit.width
}

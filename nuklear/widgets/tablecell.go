package widgets

import (
	"image/color"

	"github.com/aarzilli/nucular/label"
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
	return window.LabelWidth(label.text)
}

type CheckboxTableCell struct {
	label        string
	checked      *bool
	checkChanged func()
}

func NewCheckboxTableCell(label string, checked *bool, checkChanged func()) *CheckboxTableCell {
	return &CheckboxTableCell{
		label:        label,
		checked:      checked,
		checkChanged: checkChanged,
	}
}

func (checkbox *CheckboxTableCell) Render(window *Window) {
	if window.CheckboxText(checkbox.label, checkbox.checked) && checkbox.checkChanged != nil {
		checkbox.checkChanged()
	}
}

func (checkbox *CheckboxTableCell) MinWidth(window *Window) int {
	return window.LabelWidth(checkbox.label) + 16 // assumed width of check box
}

type LinkTableCell struct {
	text        string
	tooltipText string
	selected    *bool
	clickFunc   func(text string, window *Window)
}

func NewLinkTableCell(text, tooltipText string, clickFunc func(text string, window *Window)) *LinkTableCell {
	selected := false

	return &LinkTableCell{
		text:        text,
		tooltipText: tooltipText,
		selected:    &selected,
		clickFunc:   clickFunc,
	}
}

func (link *LinkTableCell) Render(window *Window) {
	if window.SelectableLabel(link.text, "LC", link.selected) {
		link.clickFunc(link.text, window)
	}
	if link.tooltipText != "" {
		if window.Input().Mouse.HoveringRect(window.LastWidgetBounds) {
			window.Tooltip(link.tooltipText)
		}
	}
}

func (link *LinkTableCell) MinWidth(window *Window) int {
	return window.LabelWidth(link.text)
}

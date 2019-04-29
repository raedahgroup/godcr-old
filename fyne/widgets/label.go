package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func GetLabelWithStyle(text string, align fyne.TextAlign, style fyne.TextStyle) *widget.Label {
	return widget.NewLabelWithStyle(text, align, style)
}

func (b *Box) AddLabelWithStyle(text string, align fyne.TextAlign, style fyne.TextStyle) *widget.Label {
	label := GetLabelWithStyle(text, align, style)
	b.Box.Append(label)

	return label
}

func (b *Box) AddLabel(text string) *widget.Label {
	label := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: false})
	b.Box.Append(label)

	return label
}

func (b *Box) AddBoldLabel(text string) *widget.Label {
	label := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	b.Box.Append(label)

	return label
}

func (b *Box) AddItalicLabel(text string) *widget.Label {
	label := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Italic: true, Bold: false})
	b.Box.Append(label)

	return label
}

func (b *Box) AddBoldAndItalicLabel(text string) *widget.Label {
	label := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Italic: true, Bold: true})
	b.Box.Append(label)

	return label
}

func (b *Box) AddErrorLabel(err string) *widget.Label {
	return nil
}

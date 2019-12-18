package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"image/color"
)

type Button struct {
	bar        *canvas.Rectangle
	fillColor  color.Color
	canvasText *canvas.Text

	Size      fyne.Size
	Container *ClickableBox
}

func (b *Button) MinSize() fyne.Size {
	return b.canvasText.MinSize()
}

func (b *Button) SetMinSize(size fyne.Size) {
	b.bar.SetMinSize(size)
}

func (b *Button) SetTextSize(size int) {
	b.canvasText.TextSize = size
}

func (b *Button) SetTextStyle(style fyne.TextStyle) {
	b.canvasText.TextStyle = style
}

func (b *Button) Disable() {
	b.Container.Disable()

	b.bar.FillColor = color.RGBA{196, 203, 210, 255}
}

func (b *Button) Enable() {
	b.Container.Enable()

	b.bar.FillColor = b.fillColor
	b.bar.Refresh()
}

func (b *Button) Disabled() bool {
	return b.Container.Disabled()
}

func (b *Button) SetText(text string) {
	b.canvasText.Text = text
}

// NewButton implements a custom button widget with varied size and OnTapped ability.
// To use NewButton as a fyne widget Container field should be passed.
func NewButton(fillColor color.Color, text string, OnTapped func()) *Button {
	var button Button

	button.canvasText = canvas.NewText(text, color.White)
	if fillColor == color.White {
		button.canvasText = canvas.NewText(text, color.Black)
	}
	button.canvasText.Alignment = fyne.TextAlignCenter

	button.bar = canvas.NewRectangle(fillColor)
	button.bar.SetMinSize(button.canvasText.MinSize())

	button.fillColor = fillColor

	Container := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), button.bar, button.canvasText)

	button.Container = NewClickableBox(widget.NewVBox(Container), OnTapped)

	return &button
}

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
	fillColor  color.RGBA
	canvasText *canvas.Text

	Size      fyne.Size
	Container *ClickableBox
	OnTapped  func()
}

func (b *Button) MinSize() fyne.Size {
	x := fyne.Max(b.bar.MinSize().Width, b.canvasText.MinSize().Width)
	y := fyne.Max(b.bar.MinSize().Height, b.canvasText.MinSize().Height)

	return fyne.NewSize(x, y)
}

func (b *Button) SetMinSize(size fyne.Size) {
	b.bar.SetMinSize(size)
	canvas.Refresh(b.bar)
}

func (b *Button) SetTextSize(size int) {
	b.canvasText.TextSize = size
}

func (b *Button) Disable() {
	b.Container.Disable()

	b.bar.FillColor = color.RGBA{196, 203, 210, 255}

	canvas.Refresh(b.canvasText)
	canvas.Refresh(b.bar)
	widget.Refresh(b.Container)
}

func (b *Button) Enable() {
	b.Container.Enable()

	b.bar.FillColor = b.fillColor

	canvas.Refresh(b.canvasText)
	canvas.Refresh(b.bar)
	widget.Refresh(b.Container)
}

func (b *Button) Disabled() bool {
	return b.Container.Disabled()
}

func (b *Button) SetText(text string) {
	b.canvasText.Text = text
	canvas.Refresh(b.Container)
}

func NewButton(fillColor color.RGBA, text string, OnTapped func()) *Button {
	var button Button

	button.canvasText = canvas.NewText(text, color.White)
	button.canvasText.Alignment = fyne.TextAlignCenter

	button.bar = canvas.NewRectangle(fillColor)
	button.bar.SetMinSize(button.canvasText.MinSize())

	button.fillColor = fillColor

	Container := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), button.bar, button.canvasText)

	button.Container = NewClickableBox(widget.NewVBox(Container), OnTapped)

	return &button
}

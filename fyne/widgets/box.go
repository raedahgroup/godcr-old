package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type Box struct {
	*widget.Box
	parent *Window
}

func NewHBox(children ...fyne.CanvasObject) *Box {
	return &Box{
		Box: widget.NewHBox(children...),
	}
}

func NewVBox(children ...fyne.CanvasObject) *Box {
	return &Box{
		Box: widget.NewVBox(children...),
	}
}

func NewVSpacer(height int) fyne.CanvasObject {
	space := fyne.NewSize(0, height)
	return fyne.NewContainerWithLayout(layout.NewFixedGridLayout(space), layout.NewSpacer())
}

func NewHSpacer(width int) fyne.CanvasObject {
	space := fyne.NewSize(width, 0)
	return fyne.NewContainerWithLayout(layout.NewFixedGridLayout(space), layout.NewSpacer())
}

func (b *Box) SetParent(window *Window) {
	b.parent = window
}

func (b *Box) SetTitle(title string) {
	b.AddBoldAndItalicLabel(title)
}

func (b *Box) AddButton(text string, clickFunc func()) *widget.Button {
	button := widget.NewButton(text, clickFunc)
	b.Box.Append(button)

	return button
}

func (b *Box) Add(object fyne.CanvasObject) {
	b.Box.Append(object)
}

func (b *Box) Empty() {
	b.Box.Children = []fyne.CanvasObject{}
}

func (b *Box) SetContent(box *widget.Box) {
	b.Box.Append(box)
}

func (b *Box) Update() {
	b.parent.main.Canvas().Refresh(b.Box)
}

func (b *Box) DisplayError(err error) {
	b.AddErrorLabel(err.Error())
}

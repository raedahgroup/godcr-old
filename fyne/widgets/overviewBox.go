package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
)

type overviewBox struct {
	*widget.Box

	position fyne.Position
	size     fyne.Size
}

type overviewBoxRenderer struct {
	objects []fyne.CanvasObject
	box     *overviewBox
}

func (b *overviewBox) CreateRenderer() fyne.WidgetRenderer {
	var objects []fyne.CanvasObject
	objects = append(objects, b.Box)

	return &overviewBoxRenderer{objects: objects, box: b}
}

func (r *overviewBoxRenderer) MinSize() fyne.Size {
	return fyne.NewSize(theme.Padding()*2+300, theme.Padding()*2+150)
}

func (r *overviewBoxRenderer) ApplyTheme() {

}

func (r *overviewBoxRenderer) Layout(size fyne.Size) {

}

func (r *overviewBoxRenderer) BackgroundColor() color.Color {
	return color.RGBA{242, 241, 239, 1}
}

func (r *overviewBoxRenderer) Refresh() {

}

func (r *overviewBoxRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *overviewBoxRenderer) Destroy() {
}

func (b *overviewBox) Size() fyne.Size {
	return b.size
}

func (b *overviewBox) Resize(size fyne.Size) {
	b.size = size

	if widget.Renderer(b) != nil {
		widget.Renderer(b).Layout(size)
	}
}

func (b *overviewBox) Move(pos fyne.Position) {
}

func (b *overviewBox) Show() {
}

func (b *overviewBox) Hide() {
}

func (b *overviewBox) Visible() bool {
	return true
}

func (b *overviewBox) Position() fyne.Position {
	return b.position
}

func (b *overviewBox) MinSize() fyne.Size {
	if widget.Renderer(b) == nil {
		return fyne.NewSize(0, 0)
	}
	return widget.Renderer(b).MinSize()
}

func NewOverviewBox(children ...fyne.CanvasObject) *overviewBox {
	box := &overviewBox{Box: widget.NewVBox(children...)}
	widget.Renderer(box).Layout(box.MinSize())
	return box
}

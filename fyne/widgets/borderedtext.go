package widgets

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
)

type BorderedText struct {
	Container *fyne.Container

	canvasText *canvas.Text
	border     *canvas.Rectangle
	padding    fyne.Size
}

func (borderedText *BorderedText) SetPadding(padding fyne.Size) {
	borderedText.padding = padding
	borderedText.border.SetMinSize(padding.Add(borderedText.canvasText.MinSize()))

	borderedText.border.Refresh()
	borderedText.Container.Refresh()
}

func (borderedText *BorderedText) SetText(text string) {
	borderedText.canvasText.Text = text
	borderedText.border.SetMinSize(borderedText.padding.Add(borderedText.canvasText.MinSize()))

	borderedText.canvasText.Refresh()
	borderedText.border.Refresh()
	borderedText.Container.Refresh()
}

func NewBorderedText(text string, padding fyne.Size, borderColor color.Color) (borderedText *BorderedText) {
	borderedText = &BorderedText{}
	borderedText.canvasText = canvas.NewText(text, color.White)
	borderedText.border = canvas.NewRectangle(borderColor)

	if padding == fyne.NewSize(0, 0) {
		padding = fyne.NewSize(20, 8)
	}

	borderedText.padding = padding
	borderedText.border.SetMinSize(padding.Add(borderedText.canvasText.MinSize()))
	borderedText.Container = fyne.NewContainerWithLayout(layout.NewCenterLayout(), borderedText.border, borderedText.canvasText)

	return
}

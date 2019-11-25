package widgets

import (
	"fyne.io/fyne/canvas"
	"image/color"
)

var DefaultTextColor = color.Black

func NewLargeText (text string, textColor color.Color) *canvas.Text {
	if textColor == nil {
		textColor = DefaultTextColor
	}
	return &canvas.Text{
		Color:    textColor,
		Text:     text,
		TextSize: 25,
	}
}

func NewSmallText (text string, textColor color.Color) *canvas.Text {
	if textColor == nil {
		textColor = DefaultTextColor
	}

	return &canvas.Text{
		Color:    textColor,
		Text:     text,
		TextSize: 15,
	}
}

func NewTextWithSize (text string, textColor color.Color, textSize int) *canvas.Text {
	if textColor == nil {
		textColor = DefaultTextColor
	}

	return &canvas.Text{
		Color:    textColor,
		Text:     text,
		TextSize: textSize,
	}
}

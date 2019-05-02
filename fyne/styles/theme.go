package styles

import (
	"image/color"

	"fyne.io/fyne"
	fyneTheme "fyne.io/fyne/theme"
)

type theme struct {
}

func NewTheme() fyne.Theme {
	return &theme{}
}

func (theme) BackgroundColor() color.Color {
	return DefaultBackgroundColor
}

func (theme) ButtonColor() color.Color {
	return DecredDarkBlueColor
}

func (theme) HyperlinkColor() color.Color {
	return DecredLightBlueColor
}

func (theme) TextColor() color.Color {
	return WhiteColor
}

func (theme) PlaceHolderColor() color.Color {
	return GrayColor
}

func (theme) PrimaryColor() color.Color {
	return DecredLightBlueColor
}

func (theme) FocusColor() color.Color {
	return DecredOrangeColor
}

func (theme) ScrollBarColor() color.Color {
	return DecredLightBlueColor
}

func (theme) TextSize() int {
	return 12
}

func (theme) TextFont() fyne.Resource {
	return fyneTheme.DefaultTextFont()
}

func (theme) TextBoldFont() fyne.Resource {
	return fyneTheme.DefaultTextBoldFont()
}

func (theme) TextItalicFont() fyne.Resource {
	return fyneTheme.DefaultTextBoldItalicFont()
}

func (theme) TextBoldItalicFont() fyne.Resource {
	return fyneTheme.DefaultTextBoldItalicFont()
}

func (theme) TextMonospaceFont() fyne.Resource {
	return fyneTheme.DefaultTextMonospaceFont()
}

func (theme) Padding() int {
	return 5
}

func (theme) IconInlineSize() int {
	return 20
}

func (theme) ScrollBarSize() int {
	return 10
}

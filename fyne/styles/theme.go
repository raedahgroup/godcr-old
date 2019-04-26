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
	return GrayColor
}

func (theme) ButtonColor() color.Color {
	return DecredLightBlueColor
}

func (theme) HyperlinkColor() color.Color {
	return DecredLightBlueColor
}

func (theme) TextColor() color.Color {
	return BlackColor
}

func (theme) PlaceHolderColor() color.Color {
	return GrayColor
}

func (theme) PrimaryColor() color.Color {
	return DecredDarkBlueColor
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
	return fyneTheme.DefaultTextBoldFont()
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
	return 10
}

func (theme) IconInlineSize() int {
	return 20
}

func (theme) ScrollBarSize() int {
	return 10
}

type navTheme struct {
}

func NewNavTheme() fyne.Theme {
	return &navTheme{}
}

func (navTheme) BackgroundColor() color.Color {
	return DecredDarkBlueColor
}

func (navTheme) ButtonColor() color.Color {
	return DecredLightBlueColor
}

func (navTheme) HyperlinkColor() color.Color {
	return DecredLightBlueColor
}

func (navTheme) TextColor() color.Color {
	return BlackColor
}

func (navTheme) PlaceHolderColor() color.Color {
	return GrayColor
}

func (navTheme) PrimaryColor() color.Color {
	return DecredDarkBlueColor
}

func (navTheme) FocusColor() color.Color {
	return DecredOrangeColor
}

func (navTheme) ScrollBarColor() color.Color {
	return DecredLightBlueColor
}

func (navTheme) TextSize() int {
	return 12
}

func (navTheme) TextFont() fyne.Resource {
	return fyneTheme.DefaultTextBoldFont()
}

func (navTheme) TextBoldFont() fyne.Resource {
	return fyneTheme.DefaultTextBoldFont()
}

func (navTheme) TextItalicFont() fyne.Resource {
	return fyneTheme.DefaultTextBoldItalicFont()
}

func (navTheme) TextBoldItalicFont() fyne.Resource {
	return fyneTheme.DefaultTextBoldItalicFont()
}

func (navTheme) TextMonospaceFont() fyne.Resource {
	return fyneTheme.DefaultTextMonospaceFont()
}

func (navTheme) Padding() int {
	return 10
}

func (navTheme) IconInlineSize() int {
	return 20
}

func (navTheme) ScrollBarSize() int {
	return 10
}

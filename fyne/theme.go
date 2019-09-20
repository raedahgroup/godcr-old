package fyne

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"image/color"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type appTheme struct {
	background color.Color

	button, text, icon, hyperlink, placeholder, primary, hover, scrollBar, shadow color.Color
	regular, bold, italic, bolditalic, monospace                                  fyne.Resource
	disabledButton, disabledIcon, disabledText                                    color.Color
}

func AppThem() *appTheme {
	lightThem := theme.LightTheme()
	appTheme := &appTheme{
		background:     color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		button:         color.RGBA{R: 0xd9, G: 0xd9, B: 0xd9, A: 0xff},
		disabledButton: color.RGBA{R: 0xe7, G: 0xe7, B: 0xe7, A: 0xff},
		text:           color.RGBA{R: 0x21, G: 0x21, B: 0x21, A: 0xff},
		disabledText:   color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff},
		icon:           color.RGBA{R: 0x21, G: 0x21, B: 0x21, A: 0xff},
		disabledIcon:   color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff},
		hyperlink:      color.RGBA{B: 0xd9, A: 0xff},
		placeholder:    color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
		primary:        color.RGBA{R: 41, G: 112, B: 155, A: 255},
		hover:          color.RGBA{R: 0xe7, G: 0xe7, B: 0xe7, A: 0xff},
		scrollBar:      color.RGBA{A: 0x99},
		shadow:         color.RGBA{A: 0x33},
		bold:           lightThem.TextBoldFont(),
		bolditalic:     lightThem.TextBoldItalicFont(),
		italic:         lightThem.TextBoldItalicFont(),
		monospace:      lightThem.TextMonospaceFont(),
		regular:        lightThem.TextFont(),
	}
	
	return appTheme
}

func (t *appTheme) BackgroundColor() color.Color {
	return t.background
}

// ButtonColor returns the theme's standard button colour
func (t *appTheme) ButtonColor() color.Color {
	return t.button
}

// DisabledButtonColor returns the theme's disabled button colour
func (t *appTheme) DisabledButtonColor() color.Color {
	return t.disabledButton
}

// HyperlinkColor returns the theme's standard hyperlink colour
func (t *appTheme) HyperlinkColor() color.Color {
	return t.hyperlink
}

// TextColor returns the theme's standard text colour
func (t *appTheme) TextColor() color.Color {
	return t.text
}

// DisabledIconColor returns the color for a disabledIcon UI element
func (t *appTheme) DisabledTextColor() color.Color {
	return t.disabledText
}

// IconColor returns the theme's standard text colour
func (t *appTheme) IconColor() color.Color {
	return t.icon
}

// DisabledIconColor returns the color for a disabledIcon UI element
func (t *appTheme) DisabledIconColor() color.Color {
	return t.disabledIcon
}

// PlaceHolderColor returns the theme's placeholder text colour
func (t *appTheme) PlaceHolderColor() color.Color {
	return t.placeholder
}

// PrimaryColor returns the colour used to highlight primary features
func (t *appTheme) PrimaryColor() color.Color {
	return t.primary
}

// HoverColor returns the colour used to highlight interactive elements currently under a cursor
func (t *appTheme) HoverColor() color.Color {
	return t.hover
}

// FocusColor returns the colour used to highlight a focused widget
func (t *appTheme) FocusColor() color.Color {
	return t.primary
}

// ScrollBarColor returns the color (and translucency) for a scrollBar
func (t *appTheme) ScrollBarColor() color.Color {
	return t.scrollBar
}

// ShadowColor returns the color (and translucency) for shadows used for indicating elevation
func (t *appTheme) ShadowColor() color.Color {
	return t.shadow
}

// TextSize returns the standard text size
func (t *appTheme) TextSize() int {
	return 14
}

func loadCustomFont(env, variant string, fallback fyne.Resource) fyne.Resource {
	variantPath := strings.Replace(env, "Regular", variant, 0)

	file, err := os.Open(variantPath)
	if err != nil {
		fyne.LogError("Error loading specified font", err)
		return fallback
	}
	ret, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		fyne.LogError("Error loading specified font", err2)
		return fallback
	}

	name := path.Base(variantPath)
	return &fyne.StaticResource{StaticName: name, StaticContent: ret}

}

// TextFont returns the font resource for the regular font style
func (t *appTheme) TextFont() fyne.Resource {
	return t.regular
}

// TextBoldFont retutns the font resource for the bold font style
func (t *appTheme) TextBoldFont() fyne.Resource {
	return t.bold
}

// TextItalicFont returns the font resource for the italic font style
func (t *appTheme) TextItalicFont() fyne.Resource {
	return t.italic
}

// TextBoldItalicFont returns the font resource for the bold and italic font style
func (t *appTheme) TextBoldItalicFont() fyne.Resource {
	return t.bolditalic
}

// TextMonospaceFont retutns the font resource for the monospace font face
func (t *appTheme) TextMonospaceFont() fyne.Resource {
	return t.monospace
}

// Padding is the standard gap between elements and the border around interface
// elements
func (t *appTheme) Padding() int {
	return 4
}

// IconInlineSize is the standard size of icons which appear within buttons, labels etc.
func (t *appTheme) IconInlineSize() int {
	return 20
}

// ScrollBarSize is the width (or height) of the bars on a ScrollContainer
func (t *appTheme) ScrollBarSize() int {
	return 16
}

// ScrollBarSmallSize is the width (or height) of the minimized bars on a ScrollContainer
func (t *appTheme) ScrollBarSmallSize() int {
	return 3
}

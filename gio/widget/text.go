package widget

import (
	"gioui.org/ui"
	"gioui.org/ui/layout"
	"gioui.org/ui/measure"
	"gioui.org/ui/text"
	"golang.org/x/image/font/sfnt"

	"github.com/raedahgroup/godcr/gio/helper"
)

var (
	faces measure.Faces
)

const (
	regularFontSize = 11
	bigFontSize     = 18
)

const (
	TextAlignLeft = iota
	TextAlignCenter
	TextAlignRight
)

func ErrorText(textContent string, alignment int, ctx *layout.Context) text.Label {
	return makeText(textContent, alignment, helper.Theme.DangerText, helper.Fonts.Italic, regularFontSize, ctx)
}

func SuccessText(textContent string, alignment int, ctx *layout.Context) text.Label {
	return makeText(textContent, alignment, helper.Theme.SuccessText, helper.Fonts.Italic, regularFontSize, ctx)
}

func RegularText(textContent string, alignment int, ctx *layout.Context) text.Label {
	return makeText(textContent, alignment, helper.Theme.Text, helper.Fonts.Regular, regularFontSize, ctx)
}

func ButtonText(textContent string, alignment int, ctx *layout.Context) text.Label {
	return makeText(textContent, alignment, helper.Theme.White, helper.Fonts.Regular, regularFontSize, ctx)
}

func HeadingText(textContent string, alignment int, ctx *layout.Context) text.Label {
	return makeText(textContent, alignment, helper.Theme.Text, helper.Fonts.Bold, bigFontSize, ctx)
}

func ItalicText(textContent string, alignment int, ctx *layout.Context) text.Label {
	return makeText(textContent, alignment, helper.Theme.Text, helper.Fonts.Italic, regularFontSize, ctx)
}

func RegularSecondaryText(textContent string, alignment int, ctx *layout.Context) text.Label {
	return makeText(textContent, alignment, helper.Theme.SecondaryText, helper.Fonts.Regular, regularFontSize, ctx)
}

func RegularBoldText(textContent string, alignment int, ctx *layout.Context) text.Label {
	return makeText(textContent, alignment, helper.Theme.SecondaryText, helper.Fonts.Bold, regularFontSize, ctx)
}

func RegularBoldSecondaryText(textContent string, alignment int, ctx *layout.Context) text.Label {
	return makeText(textContent, alignment, helper.Theme.SecondaryText, helper.Fonts.Regular, regularFontSize, ctx)
}

func GetTextFace(font *sfnt.Font, size float32, ctx *layout.Context) text.Face {
	faces.Reset(ctx.Config)
	return faces.For(font, ui.Sp(size))
}

func makeText(content string, alignment int, material ui.MacroOp, font *sfnt.Font, size float32, ctx *layout.Context) text.Label {
	faces.Reset(ctx)
	face := faces.For(font, ui.Sp(size))
	align := text.Start

	switch alignment {
	case TextAlignCenter:
		align = text.Middle
	case TextAlignRight:
		align = text.End
	default:
		align = text.Start
	}

	t := text.Label{
		Material:  material,
		Face:      face,
		Alignment: align,
		Text:      content,
	}
	t.Layout(ctx)

	return t
}

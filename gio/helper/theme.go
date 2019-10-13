package helper

import (
	"gioui.org/widget/material"
	"gioui.org/text"
	"gioui.org/text/opentype"
	"gioui.org/unit"

	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/goregular"
)

type Theme struct {
	*material.Theme
}

const (
	fontSize = 13
)


func NewTheme() *Theme {
	shaper := new(text.Shaper)
	shaper.Register(getFont(), opentype.Must(
		opentype.Parse(goregular.TTF),
	))
	
	shaper.Register(getItalicFont(), opentype.Must(
		opentype.Parse(goitalic.TTF),
	))

	th := material.NewTheme(shaper)
	th.Color.Primary =  DecredDarkBlueColor
	th.Color.Text =  GrayColor
	th.Color.Hint =  DecredOrangeColor
	
	th.TextSize = unit.Sp(fontSize)
	
	return &Theme{
		th,
	}
}

func getFont() text.Font {
	return text.Font{
		Size: unit.Dp(fontSize),
	}
}

func getItalicFont() text.Font {
	return text.Font {
		Size: unit.Dp(fontSize),
		Style: text.Italic,
	}
}

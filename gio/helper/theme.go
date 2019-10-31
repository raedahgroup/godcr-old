package helper

import (
	"gioui.org/text"
	"gioui.org/text/opentype"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/goregular"
)

type Theme struct {
	*material.Theme
	*text.Shaper
}

const (
	fontSize = 10
)

func NewTheme() *Theme {
	gofont.Register()

	shaper := font.Default()
	mt := material.Theme{
		Shaper: shaper,
	}

	mt.Color.Primary = DecredDarkBlueColor
	mt.Color.Text = DecredDarkBlueColor
	mt.Color.Hint = GrayColor
	mt.TextSize = unit.Px(fontSize)

	return &Theme{
		&mt,
		shaper,
	}
}

func GetFont() text.Font {
	return text.Font{
		Size: unit.Px(fontSize),
	}
}

func GetNavFont() text.Font {
	return text.Font{
		Size: unit.Dp(11),
	}
}

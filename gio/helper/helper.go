package helper 

import (
	"os"
	"image"

	"gioui.org/unit"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

const (
	logoPath = "../../gio/assets/decred.png"
	StandaloneScreenPadding = 20
)

var logo material.Image

func InitLogo(theme *Theme) error {
	logoByte, err := os.Open(logoPath)
	if err != nil {
		return err
	}

	src, _, err := image.Decode(logoByte) 
	if err != nil {
		return err
	}

	logo = theme.Image(paint.NewImageOp(src))
	logo.Scale = 1.3

	return nil
}

func DrawLogo(ctx *layout.Context) {
	inset := layout.Inset{
		Left: unit.Dp(StandaloneScreenPadding),
	}
	inset.Layout(ctx, func(){
		logo.Layout(ctx)
	})
}
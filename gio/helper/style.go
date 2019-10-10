package helper

import (
	"image/color"

	_ "image/jpeg"
	_ "image/png"

	_ "net/http/pprof"

	"gioui.org/ui"
	"gioui.org/ui/paint"

	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/sfnt"
)

var Fonts struct {
	Regular *sfnt.Font
	Bold    *sfnt.Font
	Italic  *sfnt.Font
	Mono    *sfnt.Font
}

var Theme struct {
	Text          ui.MacroOp
	SecondaryText ui.MacroOp
	Brand         ui.MacroOp
	White         ui.MacroOp
	DangerText    ui.MacroOp
	SuccessText   ui.MacroOp
}

func init() {
	Fonts.Regular = mustLoadFont(goregular.TTF)
	Fonts.Bold = mustLoadFont(gobold.TTF)
	Fonts.Italic = mustLoadFont(goitalic.TTF)
	Fonts.Mono = mustLoadFont(gomono.TTF)
	var ops ui.Ops
	Theme.Text = colorMaterial(&ops, DecredDarkBlueColor)
	Theme.SecondaryText = colorMaterial(&ops, DecredOrangeColor)
	Theme.Brand = colorMaterial(&ops, DecredDarkBlueColor)
	Theme.White = colorMaterial(&ops, WhiteColor)
	Theme.DangerText = colorMaterial(&ops, DangerColor)
	Theme.SuccessText = colorMaterial(&ops, SuccessColor)
}

func Init() {

}

func mustLoadFont(fontData []byte) *sfnt.Font {
	fnt, err := sfnt.Parse(fontData)
	if err != nil {
		panic("failed to load font")
	}
	return fnt
}

func colorMaterial(ops *ui.Ops, color color.RGBA) ui.MacroOp {
	var mat ui.MacroOp
	mat.Record(ops)
	paint.ColorOp{Color: color}.Add(ops)
	mat.Stop()
	return mat
}

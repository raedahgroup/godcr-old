package helpers

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"math"

	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	PageHeaderFont  font.Face
	PageContentFont font.Face
	NavFont         font.Face

	noPadding         = image.Point{0, 0}
	pageHeaderPadding = image.Point{20, 20}
)

const (
	scaling             = 1.9
	pageHeaderFontSize  = 13
	pageHeaderFontDPI   = 72
	pageContentFontSize = 8
	pageContentFontDPI  = 70
	navFontSize         = 11
	navFontDPI          = 62
)

func InitFonts() error {
	robotoMediumFontData, err := ioutil.ReadFile("nuklear/assets/font/Roboto-Medium.ttf")
	if err != nil {
		return err
	}

	robotoLightFontData, err := ioutil.ReadFile("nuklear/assets/font/Roboto-Light.ttf")
	if err != nil {
		return err
	}

	NavFont, err = getFont(navFontSize, navFontDPI, robotoMediumFontData)
	if err != nil {
		return err
	}

	PageHeaderFont, err = getFont(pageHeaderFontSize, pageHeaderFontDPI, robotoMediumFontData)
	if err != nil {
		return err
	}

	PageContentFont, err = getFont(pageContentFontSize, pageContentFontDPI, robotoLightFontData)
	if err != nil {
		return err
	}

	return nil
}

func getFont(fontSize, DPI int, fontData []byte) (font.Face, error) {
	ttfont, err := freetype.ParseFont(fontData)
	if err != nil {
		return nil, err
	}

	size := int(float64(fontSize) * scaling)
	options := &truetype.Options{
		Size:    float64(size),
		Hinting: font.HintingFull,
		DPI:     float64(DPI),
	}

	return truetype.NewFace(ttfont, options), nil
}

func SetFont(window *nucular.Window, font font.Face) {
	style := window.Master().Style()
	style.Font = font
	window.Master().SetStyle(style)
}

func GetStyle() *nstyle.Style {
	style := nstyle.FromTable(colorTable, scaling)

	/**window**/
	style.NormalWindow.Padding = noPadding

	/**buttons**/
	style.Button.Rounding = 0
	style.Button.Border = 0
	style.Button.TextNormal = colorWhite

	/**inputs**/
	style.Edit.Normal.Data.Color = colorWhite
	style.Edit.Active.Data.Color = colorWhite
	style.Edit.Hover.Data.Color = colorWhite
	style.Edit.Border = 1
	style.Edit.BorderColor = colorTable.ColorBorder

	/**checkbox**/
	style.Checkbox.Normal.Data.Color = colorWhite
	style.Checkbox.Active.Data.Color = colorAccent
	style.Checkbox.CursorHover.Data.Color = colorAccent
	style.Checkbox.CursorNormal.Data.Color = colorAccent
	style.Checkbox.Hover.Data.Color = colorAccent

	/**form inputs**/
	style.Edit.Border = 1
	style.Edit.Normal.Data.Color = colorWhite
	style.Edit.BorderColor = colorPrimaryBorder
	style.Edit.Active.Data.Color = colorWhite
	style.Edit.Hover.Data.Color = colorWhite
	style.Combo.Normal.Data.Color = colorWhite
	style.Combo.BorderColor = colorPrimaryBorder
	style.Combo.Active.Data.Color = colorWhite
	style.Combo.Hover.Data.Color = colorWhite

	return style
}

func SetNavStyle(window nucular.MasterWindow) {
	style := window.Style()
	// nav window background color
	style.GroupWindow.FixedBackground.Data.Color = colorNavBackground
	style.GroupWindow.Padding = noPadding

	style.Button.Padding = image.Point{33, 5}
	style.Button.Normal.Data.Color = colorPrimary
	style.Button.Hover.Data.Color = color.RGBA{7, 16, 52, 255}
	style.Button.Active.Data.Color = color.RGBA{7, 16, 52, 255}
	style.Button.TextHover = colorWhite
	style.Font = NavFont

	window.SetStyle(style)
}

func SetPageStyle(window nucular.MasterWindow) {
	style := window.Style()
	style.Button.Normal.Data.Color = colorAccent
	style.Button.Hover.Data.Color = colorAccentDark
	style.Button.Active.Data.Color = colorAccentDark
	style.GroupWindow.FixedBackground.Data.Color = colorContentBackground

	window.SetStyle(style)
}

func SetStandaloneWindowStyle(window nucular.MasterWindow) {
	style := window.Style()
	style.GroupWindow.FixedBackground.Data.Color = colorWhite
	style.GroupWindow.Padding = image.Point{20, 15}
	style.NormalWindow.ScalerSize = image.Point{50, 50}

	window.SetStyle(style)
}

func AmountToString(amount float64) string {
	amount = math.Round(amount)
	return fmt.Sprintf("%d DCR", int(amount))
}

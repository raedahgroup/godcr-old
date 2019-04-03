package helpers

import (
	"image"
	"io/ioutil"

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
	scaling             = 2.0
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
	style.Button.TextNormal = whiteColor

	/**inputs**/
	style.Edit.Normal.Data.Color = whiteColor
	style.Edit.Active.Data.Color = whiteColor
	style.Edit.Hover.Data.Color = whiteColor
	style.Edit.Border = 1
	style.Edit.BorderColor = colorTable.ColorBorder

	/**checkbox**/
	style.Checkbox.Normal.Data.Color = whiteColor
	style.Checkbox.Active.Data.Color = secondaryColor
	style.Checkbox.CursorHover.Data.Color = secondaryColor
	style.Checkbox.CursorNormal.Data.Color = secondaryColor
	style.Checkbox.Hover.Data.Color = secondaryColor

	/**form inputs**/
	style.Edit.Border = 1
	style.Edit.Normal.Data.Color = whiteColor
	style.Edit.BorderColor = borderColor
	style.Edit.Active.Data.Color = whiteColor
	style.Edit.Hover.Data.Color = whiteColor
	style.Combo.Normal.Data.Color = whiteColor
	style.Combo.BorderColor = borderColor
	style.Combo.Active.Data.Color = whiteColor
	style.Combo.Hover.Data.Color = whiteColor

	return style
}

func SetNavStyle(window nucular.MasterWindow) {
	style := window.Style()
	// nav window background color
	style.GroupWindow.FixedBackground.Data.Color = primaryColor
	style.GroupWindow.Padding = noPadding

	style.Button.Padding = image.Point{33, 2}
	style.Button.Normal.Data.Color = primaryColor
	style.Button.Hover.Data.Color = primaryColorLight
	style.Button.Active.Data.Color = primaryColorLight
	style.Button.TextHover = whiteColor
	style.Font = NavFont

	window.SetStyle(style)
}

func SetPageStyle(window nucular.MasterWindow) {
	style := window.Style()
	style.Button.Normal.Data.Color = secondaryColor
	style.Button.Hover.Data.Color = secondaryColorLight
	style.Button.Active.Data.Color = secondaryColorLight
	style.GroupWindow.FixedBackground.Data.Color = windowColor

	window.SetStyle(style)
}

func SetStandaloneWindowStyle(window nucular.MasterWindow) {
	style := window.Style()
	style.GroupWindow.FixedBackground.Data.Color = whiteColor
	style.GroupWindow.Padding = image.Point{20, 15}
	style.NormalWindow.ScalerSize = image.Point{50, 50}

	window.SetStyle(style)
}

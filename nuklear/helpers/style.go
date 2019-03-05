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
)

const (
	scaling             = 1.8
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
	style.NormalWindow.Padding = image.Point{0, 0}

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

	return style
}

func SetNavStyle(window nucular.MasterWindow) {
	style := window.Style()
	// nav window background color
	style.GroupWindow.FixedBackground.Data.Color = navBackgroundColor
	style.GroupWindow.Padding = image.Point{0, 0}

	style.Button.Padding = image.Point{33, 5}
	style.Button.Hover.Data.Color = color.RGBA{7, 16, 52, 255}
	style.Button.Active.Data.Color = color.RGBA{7, 16, 52, 255}
	style.Button.TextHover = whiteColor
	style.Font = NavFont

	window.SetStyle(style)
}

func SetPageStyle(window nucular.MasterWindow) {
	style := window.Style()
	style.GroupWindow.FixedBackground.Data.Color = contentBackgroundColor

	window.SetStyle(style)
}

func SetStandaloneWindowStyle(window nucular.MasterWindow) {
	style := window.Style()
	style.GroupWindow.FixedBackground.Data.Color = whiteColor
	style.GroupWindow.Padding = image.Point{20, 15}
	style.NormalWindow.ScalerSize = image.Point{50, 50}

	window.SetStyle(style)
}

func AmountToString(amount float64) string {
	amount = math.Round(amount)
	return fmt.Sprintf("%d DCR", int(amount))
}

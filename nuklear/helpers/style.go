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
	whiteColor             = color.RGBA{0xff, 0xff, 0xff, 0xff}
	navBackgroundColor     = color.RGBA{9, 20, 64, 255}
	contentBackgroundColor = color.RGBA{240, 240, 250, 255}
	PageHeaderFont         font.Face
	PageContentFont        font.Face
	NavFont                font.Face
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

var colorTable = nstyle.ColorTable{
	ColorText:                  color.RGBA{106, 106, 106, 255},
	ColorWindow:                contentBackgroundColor,
	ColorHeader:                color.RGBA{175, 175, 175, 255},
	ColorBorder:                color.RGBA{0, 0, 0, 255},
	ColorButton:                color.RGBA{9, 20, 64, 255},
	ColorButtonHover:           color.RGBA{255, 255, 255, 255},
	ColorButtonActive:          color.RGBA{0, 153, 204, 255},
	ColorToggle:                color.RGBA{150, 150, 150, 255},
	ColorToggleHover:           color.RGBA{120, 120, 120, 255},
	ColorToggleCursor:          color.RGBA{175, 175, 175, 255},
	ColorSelect:                color.RGBA{175, 175, 175, 255},
	ColorSelectActive:          color.RGBA{190, 190, 190, 255},
	ColorSlider:                color.RGBA{190, 190, 190, 255},
	ColorSliderCursor:          color.RGBA{80, 80, 80, 255},
	ColorSliderCursorHover:     color.RGBA{70, 70, 70, 255},
	ColorSliderCursorActive:    color.RGBA{60, 60, 60, 255},
	ColorProperty:              color.RGBA{175, 175, 175, 255},
	ColorEdit:                  color.RGBA{150, 150, 150, 255},
	ColorEditCursor:            color.RGBA{0, 0, 0, 255},
	ColorCombo:                 color.RGBA{175, 175, 175, 255},
	ColorChart:                 color.RGBA{160, 160, 160, 255},
	ColorChartColor:            color.RGBA{45, 45, 45, 255},
	ColorChartColorHighlight:   color.RGBA{255, 0, 0, 255},
	ColorScrollbar:             color.RGBA{180, 180, 180, 255},
	ColorScrollbarCursor:       color.RGBA{140, 140, 140, 255},
	ColorScrollbarCursorHover:  color.RGBA{150, 150, 150, 255},
	ColorScrollbarCursorActive: color.RGBA{160, 160, 160, 255},
	ColorTabHeader:             color.RGBA{0x89, 0x89, 0x89, 0xff},
}

func InitFonts() error {
	robotoMediumFontData, err := readFontFile("nuklear/assets/font/Roboto-Medium.ttf")
	if err != nil {
		return err
	}

	robotoLightFontData, err := readFontFile("nuklear/assets/font/Roboto-Light.ttf")
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

func readFontFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
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

func SetPageStyle(w nucular.MasterWindow) {
	style := w.Style()
	style.GroupWindow.FixedBackground.Data.Color = contentBackgroundColor

	w.SetStyle(style)
}

func AmountToString(amount float64) string {
	amount = math.Round(amount)
	return fmt.Sprintf("%d DCR", int(amount))
}

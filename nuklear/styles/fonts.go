package styles

import (
	"io/ioutil"

	"github.com/aarzilli/nucular"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	NavFont         font.Face
	PageHeaderFont  font.Face
	PageContentFont font.Face
)

const (
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

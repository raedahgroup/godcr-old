package styles

import (
	"io/ioutil"

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
	pageHeaderFontSize  = 18
	pageContentFontSize = 14
	navFontSize         = 14
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

	NavFont, err = getFont(navFontSize, robotoMediumFontData)
	if err != nil {
		return err
	}

	PageHeaderFont, err = getFont(pageHeaderFontSize, robotoMediumFontData)
	if err != nil {
		return err
	}

	PageContentFont, err = getFont(pageContentFontSize, robotoLightFontData)
	if err != nil {
		return err
	}

	return nil
}

func getFont(fontSize int, fontData []byte) (font.Face, error) {
	ttfont, err := freetype.ParseFont(fontData)
	if err != nil {
		return nil, err
	}

	options := &truetype.Options{
		Size: float64(fontSize),
	}

	return truetype.NewFace(ttfont, options), nil
}

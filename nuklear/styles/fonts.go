package styles

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/raedahgroup/godcr/nuklear/assets"
	"golang.org/x/image/font"
)

var (
	NavFont             font.Face
	PageHeaderFont      font.Face
	PageContentFont     font.Face
	BoldPageContentFont font.Face
)

const (
	pageHeaderFontSize  = 18
	pageContentFontSize = 16
	navFontSize         = 16
)

func InitFonts() error {
	fontsBytes, err := assets.GetFonts(assets.SourceSansProRegular, assets.SourceSansProSemibold, assets.SourceSansProSemiboldIt)
	if err != nil {
		return err
	}

	regularFontBytes := fontsBytes[assets.SourceSansProRegular]
	semiBoldFontBytes := fontsBytes[assets.SourceSansProSemibold]
	boldItalicsFontBytes := fontsBytes[assets.SourceSansProSemiboldIt]

	NavFont, err = getFont(navFontSize, regularFontBytes)
	if err != nil {
		return err
	}

	PageHeaderFont, err = getFont(pageHeaderFontSize, boldItalicsFontBytes)
	if err != nil {
		return err
	}

	PageContentFont, err = getFont(pageContentFontSize, regularFontBytes)
	if err != nil {
		return err
	}

	BoldPageContentFont, err = getFont(pageContentFontSize, semiBoldFontBytes)
	if err != nil {
		return err
	}

	return nil
}

func getFont(fontSize int, fontBytes []byte) (font.Face, error) {
	ttfont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	options := &truetype.Options{
		Size: float64(fontSize),
	}

	return truetype.NewFace(ttfont, options), nil
}

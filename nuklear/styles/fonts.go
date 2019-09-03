package styles

import (
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	NavFont                  font.Face
	PageHeaderFont           font.Face
	PageContentFont          font.Face
	BoldPageContentFont      font.Face
	SmallBoldPageContentFont font.Face
	LightPageContentFont     font.Face
)

const (
	pageHeaderFontSize  = 18
	pageContentFontSize = 16
	navFontSize         = 16
	SmallFontSize       = 13
)

// todo fix font file paths
func InitFonts() error {
	boldItalicsFontBytes, err := ioutil.ReadFile("../../nuklear/assets/font/SourceSansPro-SemiboldIt.ttf")
	if err != nil {
		return err
	}

	semiBoldFontBytes, err := ioutil.ReadFile("../../nuklear/assets/font/SourceSansPro-Semibold.ttf")
	if err != nil {
		return err
	}

	regularFontBytes, err := ioutil.ReadFile("../../nuklear/assets/font/SourceSansPro-Regular.ttf")
	if err != nil {
		return err
	}

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

	SmallBoldPageContentFont, err = getFont(SmallFontSize, semiBoldFontBytes)
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

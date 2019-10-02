package assets

import "github.com/gobuffalo/packr/v2"

const (
	SourceSansProSemiboldIt = "SourceSansPro-SemiboldIt.ttf"
	SourceSansProSemibold   = "SourceSansPro-Semibold.ttf"
	SourceSansProRegular    = "SourceSansPro-Regular.ttf"
)

var fontsBox = packr.New("fonts", "fonts")

// GetFonts returns a map from the names of the fonts passed as arguments to
// the font resources that correspond to them. If an error is encountered
// while loading any of the fonts, the error is returned immediately.
func GetFonts(names ...string) (map[string][]byte, error) {
	fonts := make(map[string][]byte, len(names))
	for _, name := range names {
		fontBytes, err := fontsBox.Find(name)
		if err != nil {
			return nil, err
		}
		fonts[name] = fontBytes
	}
	return fonts, nil
}

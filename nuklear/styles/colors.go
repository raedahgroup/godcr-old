package styles

import (
	"image/color"

	"github.com/aarzilli/nucular/style"
)

var (
	WhiteColor         = color.RGBA{255, 255, 255, 255}
	BlackColor         = color.RGBA{0, 0, 0, 255}
	GrayColor          = color.RGBA{128, 128, 128, 255}
	AlternateGrayColor = color.RGBA{248, 249, 250, 255}

	BorderColor = color.RGBA{0, 0, 0, 25}

	DecredDarkBlueColor  = color.RGBA{9, 20, 64, 255}
	DecredLightBlueColor = color.RGBA{112, 203, 255, 255}

	DecredOrangeColor = color.RGBA{237, 109, 71, 255}
	DecredGreenColor  = color.RGBA{65, 191, 83, 255}
)

// masterWindowColorTable describes default colors for various widgets
var masterWindowColorTable = style.ColorTable{
	// background, texts and borders
	ColorWindow: WhiteColor,
	ColorText:   BlackColor,
	ColorBorder: GrayColor,

	// edits (input fields)
	ColorEdit:       WhiteColor,
	ColorEditCursor: BlackColor,

	// toggles (checkboxes)
	ColorToggle:       GrayColor,
	ColorToggleHover:  DecredDarkBlueColor,
	ColorToggleCursor: DecredDarkBlueColor,

	// combo (dropdowns)
	ColorCombo: WhiteColor,

	// progress bar
	ColorSlider:       GrayColor,
	ColorSliderCursor: DecredLightBlueColor,
}

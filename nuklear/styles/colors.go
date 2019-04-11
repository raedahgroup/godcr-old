package styles

import (
	"image/color"

	"github.com/aarzilli/nucular/style"
)

var (
	WhiteColor = color.RGBA{255, 255, 255, 255}
	BlackColor = color.RGBA{0, 0, 0, 255}
	GrayColor  = color.RGBA{200, 200, 200, 255}

	DecredDarkBlueColor = color.RGBA{9, 20, 64, 255} // used in nav background and button background
	//SemiTransparentDecredDarkBlueColor = color.RGBA{9, 20, 60, 242}     // used for button background on hover
	DecredLightBlueColor = color.RGBA{112, 203, 255, 255} // used for border and for checkbox checked state

	DecredOrangeColor = color.RGBA{237, 109, 71, 255}
	DecredGreenColor  = color.RGBA{65, 191, 83, 255}
)

// masterWindowColorTable describes default colors for various widgets
var masterWindowColorTable = style.ColorTable{
	// background, texts and borders
	ColorWindow: WhiteColor,
	ColorText:   BlackColor,
	ColorBorder: DecredLightBlueColor,

	// edits (input fields)
	ColorEdit:       WhiteColor,
	ColorEditCursor: BlackColor,

	// toggles (checkboxes)
	ColorToggle:       WhiteColor,
	ColorToggleHover:  DecredLightBlueColor,
	ColorToggleCursor: DecredLightBlueColor,

	// combo (dropdowns)
	ColorCombo: WhiteColor,

	// progress bar
	ColorSlider:       GrayColor,
	ColorSliderCursor: DecredLightBlueColor,
}

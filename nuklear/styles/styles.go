package styles

import (
	"image"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/style"
)

var (
	noPadding         = image.Point{0, 0}
	buttonPadding = image.Point{10, 10}
)

const (
	scaling             = 2.0
)

func MasterWindowStyle() *style.Style {
	// load colors from color table then set other styles
	masterWindowStyle := style.FromTheme()

	// style windows
	masterWindowStyle.NormalWindow.Padding = noPadding
	masterWindowStyle.GroupWindow.Padding = noPadding
	masterWindowStyle.TooltipWindow.Padding = noPadding
	masterWindowStyle.ComboWindow.Padding = noPadding

	// style buttons
	masterWindowStyle.Button.Rounding = 0
	masterWindowStyle.Button.Border = 0
	masterWindowStyle.Button.Padding = buttonPadding
	masterWindowStyle.Button.TextNormal = WhiteColor // button should not use default text color
	masterWindowStyle.Button.TextHover = WhiteColor
	masterWindowStyle.Button.TextActive = WhiteColor

	// style input fields
	masterWindowStyle.Edit.Border = 1

	return masterWindowStyle
}

func SetNavStyle(masterWindow nucular.MasterWindow) {
	currentStyle := masterWindow.Style()

	// style the group window that will hold the nav buttons
	currentStyle.GroupWindow.FixedBackground.Data.Color = DecredDarkBlueColor

	// set the font for the nav buttons
	currentStyle.Font = NavFont

	// apply the updated style
	masterWindow.SetStyle(currentStyle)
}

func SetPageStyle(masterWindow nucular.MasterWindow) {
	currentStyle := masterWindow.Style()
	currentStyle.GroupWindow.FixedBackground.Data.Color = WhiteColor // todo this needed?
	masterWindow.SetStyle(currentStyle)
}

// todo prolly delete this function
func SetStandaloneWindowStyle(window nucular.MasterWindow) {
	standAloneWindowStyle := window.Style()

	standAloneWindowStyle.GroupWindow.Padding = noPadding
	//standAloneWindowStyle.NormalWindow.ScalerSize = image.Point{50, 50}

	window.SetStyle(standAloneWindowStyle)
}

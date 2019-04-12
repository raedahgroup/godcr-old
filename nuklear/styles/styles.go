package styles

import (
	"image"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/style"
)

var (
	noPadding     = image.Point{0, 0}
	buttonPadding = image.Point{10, 10}
)

func MasterWindowStyle() *style.Style {
	// load colors from color table then set other styles
	masterWindowStyle := style.FromTable(masterWindowColorTable, 1.0)

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

	// style progress bars
	masterWindowStyle.Progress.Padding = noPadding

	// style checkbox
	masterWindowStyle.Checkbox.Padding = noPadding

	return masterWindowStyle
}

func SetNavStyle(masterWindow nucular.MasterWindow) {
	currentStyle := masterWindow.Style()
	currentStyle.Font = NavFont
	currentStyle.GroupWindow.Padding = image.Point{0, 0}
	currentStyle.GroupWindow.FixedBackground.Data.Color = DecredDarkBlueColor
	currentStyle.Button.Normal.Data.Color = DecredDarkBlueColor
	currentStyle.Button.Hover.Data.Color = DecredLightBlueColor
	currentStyle.Button.Active.Data.Color = DecredLightBlueColor
}

func SetPageStyle(masterWindow nucular.MasterWindow) {
	currentStyle := masterWindow.Style()
	currentStyle.GroupWindow.FixedBackground.Data.Color = WhiteColor
	currentStyle.Button.Normal.Data.Color = DecredLightBlueColor
	currentStyle.Button.Hover.Data.Color = DecredDarkBlueColor
	currentStyle.Button.Active.Data.Color = DecredDarkBlueColor
}

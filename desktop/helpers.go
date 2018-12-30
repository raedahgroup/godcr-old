package desktop

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
	nstyle "github.com/aarzilli/nucular/style"
	"golang.org/x/image/font"
)

var (
	whiteColor             = color.RGBA{0xff, 0xff, 0xff, 0xff}
	navBackgroundColor     = color.RGBA{9, 20, 64, 255}
	contentBackgroundColor = color.RGBA{240, 240, 250, 255}
	fontSize               = 13
	defaultFont            font.Face
)

const (
	scaling = 1.8
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

type window struct {
	*nucular.Window
}

func newWindow(title string, w *nucular.Window, flags nucular.WindowFlags) *window {
	if nw := w.GroupBegin(title, flags); nw != nil {
		return &window{
			nw,
		}
	}
	return nil
}

func (w *window) header(title string) {
	w.Row(40).Dynamic(1)
	bounds := rect.Rect{
		X: contentArea.X,
		Y: contentArea.Y,
		W: contentArea.W,
		H: 80,
	}

	_, out := w.Custom(nstyle.WidgetStateActive)
	if out != nil {
		out.FillRect(bounds, 0, whiteColor)
	}

	font := w.Master().Style().Font
	bounds.Y += 25
	bounds.X += 30

	out.DrawText(bounds, title, font, colorTable.ColorText)
}

func (w *window) contentWindow(title string) *window {
	w.Row(0).Dynamic(1)
	w.style()
	return newWindow(title, w.Window, 0)
}

func (w *window) setErrorMessage(message string) {
	w.Row(300).Dynamic(1)
	w.LabelWrap(message)
}

func (w *window) style() {
	style := w.Master().Style()
	style.GroupWindow.Padding = image.Point{20, 20}

	w.Master().SetStyle(style)
}

func (w *window) end() {
	w.GroupEnd()
}

func getStyle() *nstyle.Style {
	style := nstyle.FromTable(colorTable, scaling)

	/**window**/
	style.NormalWindow.Padding = image.Point{0, 0}

	/**buttons**/
	style.Button.Rounding = 0
	style.Button.Border = 0
	style.Button.TextNormal = whiteColor

	return style
}

func setNavStyle(window nucular.MasterWindow) {
	style := window.Style()
	// nav window background color
	style.GroupWindow.FixedBackground.Data.Color = navBackgroundColor
	style.GroupWindow.Padding = image.Point{0, 0}

	style.Button.Padding = image.Point{43, 20}
	style.Button.Hover.Data.Color = color.RGBA{7, 16, 52, 255}
	style.Button.Active.Data.Color = color.RGBA{7, 16, 52, 255}
	style.Button.TextHover = whiteColor

	window.SetStyle(style)
}

func (d *Desktop) setPageStyle() {
	style := d.window.Style()
	style.GroupWindow.FixedBackground.Data.Color = contentBackgroundColor

	d.window.SetStyle(style)
}

func amountToString(amount float64) string {
	amount = math.Round(amount)
	return fmt.Sprintf("%d DCR", int(amount))
}

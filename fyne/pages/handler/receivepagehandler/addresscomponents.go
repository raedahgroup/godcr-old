package receivepagehandler

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (receivePage *ReceivePageObjects) initQrImageAndAddress() {
	receivePage.qrImage = widget.NewIcon(theme.FyneLogo())
	receivePage.address = widgets.NewTextWithStyle("", color.RGBA{41, 112, 255, 255}, fyne.TextStyle{Bold: true}, fyne.TextAlignCenter, 15)

	receivePage.ReceivePageContents.Append(widget.NewHBox(widgets.NewHSpacer(50), layout.NewSpacer(),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(300, 300)), receivePage.qrImage), layout.NewSpacer()))
	receivePage.ReceivePageContents.Append(widgets.NewVSpacer(10))
	receivePage.ReceivePageContents.Append(widget.NewHBox(widgets.NewHSpacer(50), receivePage.address))

	receivePage.generateAddressAndQR(false)
}

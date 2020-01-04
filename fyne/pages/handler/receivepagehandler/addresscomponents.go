package receivepagehandler

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (receivePage *ReceivePageObjects) initQrImageAndAddress() {
	receivePage.qrImage = widget.NewIcon(theme.FyneLogo())
	receivePage.address = widgets.NewTextWithStyle("", values.Blue, fyne.TextStyle{Bold: true}, fyne.TextAlignCenter, values.DefaultTextSize)

	receivePage.borderedContent.Append(widgets.NewHBox(layout.NewSpacer(),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(300, 300)), receivePage.qrImage), layout.NewSpacer()))
	receivePage.borderedContent.Append(widgets.NewVSpacer(values.SpacerSize10))
	receivePage.borderedContent.Append(receivePage.address)

	receivePage.generateAddressAndQR(false)
}

package receivepagehandler

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (receivePage *ReceivePageObjects) initTapToCopyText() {
	tapToCopy := widgets.NewClickableWidget(widget.NewLabelWithStyle(values.TapToCopy, fyne.TextAlignCenter, fyne.TextStyle{}), func() {
		clipboard := receivePage.Window.Clipboard()
		clipboard.SetContent(receivePage.address.Text)

		receivePage.showInfoLabel(receivePage.addressCopiedLabel, "Address copied")
	})

	receivePage.borderedContent.Append(widgets.NewHBox(layout.NewSpacer(), tapToCopy, layout.NewSpacer()))
}

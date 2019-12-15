package receivepagehandler

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (receivePage *ReceivePageObjects) initTapToCopyText() {
	tapToCopy := widgets.NewClickableBox(widget.NewHBox(widget.NewLabelWithStyle(values.TapToCopy, fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})), func() {
		clipboard := receivePage.Window.Clipboard()
		clipboard.SetContent(receivePage.address.Text)

		receivePage.showInfoLabel(receivePage.addressCopiedLabel, "Address copied")
	})

	receivePage.ReceivePageContents.Append(widget.NewHBox(widgets.NewHSpacer(50), layout.NewSpacer(), tapToCopy, layout.NewSpacer()))
}

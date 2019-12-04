package sendpagehandler

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func BaseWidgets(infoIcon, moreIcon *fyne.StaticResource, amountEntry, destinationAddressEntry *widget.Entry, window fyne.Window) *widget.Box {
	// define base widget consisting of label, more icon and info button
	sendLabel := widget.NewLabelWithStyle("Send DCR", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true, Italic: true})

	dialogLabel := widget.NewLabelWithStyle("Input the destination \nwallet address and the amount in \nDCR to send funds.", fyne.TextAlignLeading, fyne.TextStyle{})

	var clickabelInfoIcon *widgets.ImageButton
	clickabelInfoIcon = widgets.NewImageButton(infoIcon, nil, func() {
		var popup *widget.PopUp
		confirmationText := canvas.NewText("Got it", color.RGBA{41, 112, 255, 255})
		confirmationText.TextStyle.Bold = true

		dialog := widget.NewVBox(
			widgets.NewVSpacer(12),
			widget.NewLabelWithStyle("Send DCR", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widgets.NewVSpacer(30),
			dialogLabel,
			widget.NewHBox(layout.NewSpacer(), widgets.NewClickableBox(widget.NewVBox(confirmationText), func() { popup.Hide() })),
			widgets.NewVSpacer(10))

		popup = widget.NewModalPopUp(widget.NewHBox(widgets.NewHSpacer(24), dialog, widgets.NewHSpacer(20)), window.Canvas())
	})

	var clickabelMoreIcon *widgets.ImageButton
	clickabelMoreIcon = widgets.NewImageButton(moreIcon, nil, func() {
		var popup *widget.PopUp
		popup = widget.NewPopUp(widgets.NewButton(color.White, "Clear all fields", func() {
			amountEntry.SetText("")
			destinationAddressEntry.SetText("")
			popup.Hide()

		}).Container, window.Canvas())
		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			clickabelMoreIcon).Add(fyne.NewPos(clickabelMoreIcon.MinSize().Width, clickabelMoreIcon.MinSize().Height)))
	})

	return widget.NewHBox(sendLabel, layout.NewSpacer(), clickabelInfoIcon, widgets.NewHSpacer(26), clickabelMoreIcon)
}

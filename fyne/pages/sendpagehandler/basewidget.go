package sendpagehandler

import (
	"errors"
	"image/color"

	"github.com/raedahgroup/godcr/fyne/assets"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initBaseObjects() error {
	icons, err := assets.GetIcons(assets.InfoIcon, assets.MoreIcon)
	if err != nil {
		return errors.New("Could not load base object icons")
	}
	// define base widget consisting of label, more icon and info button
	sendLabel := widget.NewLabelWithStyle(sendDcr, fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	dialogLabel := widget.NewLabelWithStyle(sendPageInfo, fyne.TextAlignLeading, fyne.TextStyle{})

	var clickabelInfoIcon *widgets.ImageButton
	clickabelInfoIcon = widgets.NewImageButton(icons[assets.InfoIcon], nil, func() {
		var popup *widget.PopUp
		confirmationText := canvas.NewText("Got it", color.RGBA{41, 112, 255, 255})
		confirmationText.TextStyle.Bold = true

		dialog := widget.NewVBox(
			widgets.NewVSpacer(12),
			widget.NewLabelWithStyle(sendDcr, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widgets.NewVSpacer(30),
			dialogLabel,
			widget.NewHBox(layout.NewSpacer(), widgets.NewClickableBox(widget.NewVBox(confirmationText), func() { popup.Hide() })),
			widgets.NewVSpacer(10))

		popup = widget.NewModalPopUp(widget.NewHBox(widgets.NewHSpacer(24), dialog, widgets.NewHSpacer(20)), sendPage.Window.Canvas())
	})

	var clickabelMoreIcon *widgets.ImageButton
	clickabelMoreIcon = widgets.NewImageButton(icons[assets.MoreIcon], nil, func() {
		var popup *widget.PopUp
		popup = widget.NewPopUp(widgets.NewButton(color.White, "Clear all fields", func() {
			sendPage.amountEntry.SetText("")
			sendPage.destinationAddressEntry.SetText("")
			popup.Hide()

		}).Container, sendPage.Window.Canvas())
		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			clickabelMoreIcon).Add(fyne.NewPos(clickabelMoreIcon.MinSize().Width, clickabelMoreIcon.MinSize().Height)))
	})

	sendPage.SendPageContents.Append(widget.NewHBox(sendLabel, layout.NewSpacer(), clickabelInfoIcon, widgets.NewHSpacer(26), clickabelMoreIcon))
	return err
}

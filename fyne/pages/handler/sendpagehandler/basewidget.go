package sendpagehandler

import (
	"errors"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/constantvalues"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initBaseObjects() error {
	icons, err := assets.GetIcons(assets.InfoIcon, assets.MoreIcon)
	if err != nil {
		return errors.New(constantvalues.BaseObjectsIconErr)
	}
	// define base widget consisting of label, more icon and info button
	sendLabel := widget.NewLabelWithStyle(constantvalues.SendDcr, fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	var clickableInfoIcon *widgets.ImageButton
	clickableInfoIcon = widgets.NewImageButton(icons[assets.InfoIcon], nil, func() {
		var popup *widget.PopUp

		dialogLabel := widget.NewLabelWithStyle(constantvalues.SendPageInfo, fyne.TextAlignLeading, fyne.TextStyle{})
		confirmationText := canvas.NewText(constantvalues.GotIt, color.RGBA{41, 112, 255, 255})
		confirmationText.TextStyle.Bold = true

		dialog := widget.NewVBox(
			widgets.NewVSpacer(12),
			widget.NewLabelWithStyle(constantvalues.SendDcr, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widgets.NewVSpacer(30),
			dialogLabel,
			widget.NewHBox(layout.NewSpacer(), widgets.NewClickableBox(widget.NewVBox(confirmationText), func() { popup.Hide() })),
			widgets.NewVSpacer(10))

		popup = widget.NewModalPopUp(widget.NewHBox(widgets.NewHSpacer(24), dialog, widgets.NewHSpacer(20)), sendPage.Window.Canvas())
	})

	var clickableMoreIcon *widgets.ImageButton
	clickableMoreIcon = widgets.NewImageButton(icons[assets.MoreIcon], nil, func() {
		var popup *widget.PopUp
		popup = widget.NewPopUpAtPosition(
			widgets.NewButton(color.White, constantvalues.ClearField, func() {
				sendPage.amountEntry.SetText("")
				sendPage.destinationAddressEntry.SetText("")
				popup.Hide()

			}).Container, sendPage.Window.Canvas(), fyne.CurrentApp().Driver().AbsolutePositionForObject(
				clickableMoreIcon).Add(fyne.NewPos(clickableMoreIcon.MinSize().Width, clickableMoreIcon.MinSize().Height)))

	})

	sendPage.SendPageContents.Append(widget.NewHBox(sendLabel, layout.NewSpacer(), clickableInfoIcon, widgets.NewHSpacer(15), clickableMoreIcon))
	return err
}

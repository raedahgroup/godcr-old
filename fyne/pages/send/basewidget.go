package send

import (
	"errors"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initBaseObjects() error {
	icons, err := assets.GetIcons(assets.InfoIcon, assets.MoreIcon)
	if err != nil {
		return errors.New(values.BaseObjectsIconErr)
	}
	// define base widget consisting of label, more icon and info button
	sendLabel := widget.NewLabelWithStyle(values.SendDcr, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	var clickableInfoIcon *widgets.ImageButton
	clickableInfoIcon = widgets.NewImageButton(icons[assets.InfoIcon], nil, func() {
		var popup *widget.PopUp

		dialogLabel := widget.NewLabelWithStyle(values.SendPageInfo, fyne.TextAlignLeading, fyne.TextStyle{})
		confirmationText := widgets.NewTextWithStyle(values.GotIt, values.Blue, fyne.TextStyle{Bold: true}, fyne.TextAlignLeading, values.DefaultTextSize)

		dialog := widget.NewVBox(
			widgets.NewVSpacer(values.SpacerSize10),
			widget.NewLabelWithStyle(values.SendDcr, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widgets.NewVSpacer(values.SpacerSize10),
			dialogLabel,
			widgets.NewVSpacer(values.SpacerSize10),
			widget.NewHBox(layout.NewSpacer(), widgets.NewClickableBox(widget.NewHBox(confirmationText), func() { popup.Hide() })),
			widgets.NewVSpacer(values.SpacerSize10))

		popup = widget.NewModalPopUp(widget.NewHBox(widgets.NewHSpacer(values.SpacerSize20), dialog, widgets.NewHSpacer(values.SpacerSize20)), sendPage.Window.Canvas())
	})

	var clickableMoreIcon *widgets.ImageButton
	clickableMoreIcon = widgets.NewImageButton(icons[assets.MoreIcon], nil, func() {
		var popup *widget.PopUp
		popup = widget.NewPopUp(
			widgets.NewClickableBox(widget.NewHBox(widget.NewLabel(values.ClearField)), func() {
				sendPage.amountEntry.SetText("")
				sendPage.destinationAddressEntry.SetText("")
				popup.Hide()

			}), sendPage.Window.Canvas())

		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			clickableMoreIcon).Add(fyne.NewPos(10, clickableMoreIcon.MinSize().Height+5).Subtract(fyne.NewPos(popup.MinSize().Width, 0))))
	})

	sendPage.SendPageContents.Append(widget.NewHBox(sendLabel, layout.NewSpacer(), clickableInfoIcon, widgets.NewHSpacer(values.SpacerSize16), clickableMoreIcon))
	return err
}

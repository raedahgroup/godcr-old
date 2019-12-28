package receive

import (
	"errors"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (receivePage *ReceivePageObjects) initBaseObjects() error {
	icons, err := assets.GetIcons(assets.InfoIcon, assets.MoreIcon)
	if err != nil {
		return errors.New(values.BaseObjectsIconErr)
	}

	receivePageLabel := widget.NewLabelWithStyle(values.ReceivePageLabel, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	clickableInfoIcon := widgets.NewImageButton(icons[assets.InfoIcon], nil, func() {
		var popup *widget.PopUp

		dialogLabel := widget.NewLabelWithStyle(values.ReceivePageInfo, fyne.TextAlignLeading, fyne.TextStyle{})
		confirmationText := canvas.NewText(values.GotIt, values.Blue)
		confirmationText.TextStyle.Bold = true

		dialog := widget.NewVBox(
			widgets.NewVSpacer(values.SpacerSize10),
			widget.NewLabelWithStyle(values.ReceivePageLabel, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widgets.NewVSpacer(values.SpacerSize10),
			dialogLabel,
			widgets.NewVSpacer(values.SpacerSize10),
			widget.NewHBox(layout.NewSpacer(), widgets.NewClickableBox(widget.NewHBox(confirmationText), func() { popup.Hide() })),
			widgets.NewVSpacer(values.SpacerSize10))

		popup = widget.NewModalPopUp(widget.NewHBox(widgets.NewHSpacer(values.SpacerSize24), dialog, widgets.NewHSpacer(values.SpacerSize20)), receivePage.Window.Canvas())
	})

	var clickableMoreIcon *widgets.ImageButton
	clickableMoreIcon = widgets.NewImageButton(icons[assets.MoreIcon], nil, func() {
		var popup *widget.PopUp

		popup = widget.NewPopUp(
			widgets.NewClickableBox(widget.NewHBox(widget.NewLabel(values.GenerateNewAddress)), func() {
				receivePage.generateAddressAndQR(true)
				popup.Hide()

			}), receivePage.Window.Canvas())
		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			clickableMoreIcon).Add(fyne.NewPos(10, clickableMoreIcon.MinSize().Height+5).Subtract(fyne.NewPos(popup.MinSize().Width, 0))))
	})

	receivePage.ReceivePageContents.Append(widget.NewHBox(receivePageLabel, layout.NewSpacer(), clickableInfoIcon, widgets.NewHSpacer(values.SpacerSize16), clickableMoreIcon))
	return err
}

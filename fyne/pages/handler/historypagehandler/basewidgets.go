package historypagehandler

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (historyPage *HistoryPageData) initBaseObjects() error {
	icons, err := assets.GetIcons(assets.CollapseIcon, assets.InfoIcon, assets.SendIcon, assets.ReceiveIcon, assets.ReceiveIcon, assets.InfoIcon, assets.RedirectIcon)
	if err != nil {
		return err
	}
	historyPage.icons = icons
	// history page title label
	historyTitleLabel := widget.NewLabelWithStyle(values.HistoryTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})

	// infoPopUp creates a popup with history page hint-text
	var infoIcon *widgets.ImageButton
	infoIcon = widgets.NewImageButton(historyPage.icons[assets.InfoIcon], nil, func() {
		infoLabel := canvas.NewText(values.TxdetailsHint, values.DefaultTextColor)
		infoLabel.TextSize = 12
		infoLabel.TextStyle = fyne.TextStyle{Monospace: true}

		info2Label := canvas.NewText(values.CopyHint, values.DefaultTextColor)
		info2Label.TextSize = 12
		info2Label.TextStyle = fyne.TextStyle{Monospace: true}

		gotItLabel := canvas.NewText(values.GotIt, values.Blue)
		gotItLabel.TextStyle = fyne.TextStyle{Bold: true}
		gotItLabel.TextSize = 12

		var infoPopUp *widget.PopUp
		infoPopUp = widget.NewPopUp(widget.NewVBox(
			widgets.NewVSpacer(values.SpacerSize2),
			widget.NewHBox(widgets.NewHSpacer(values.SpacerSize2), infoLabel, widgets.NewHSpacer(values.SpacerSize2)),
			widgets.NewVSpacer(values.SpacerSize2),
			widget.NewHBox(widgets.NewHSpacer(values.SpacerSize2), info2Label, widgets.NewHSpacer(values.SpacerSize2)),
			widgets.NewVSpacer(values.SpacerSize2),
			widget.NewHBox(layout.NewSpacer(), widgets.NewClickableBox(widget.NewVBox(gotItLabel), func() { infoPopUp.Hide() }), widgets.NewHSpacer(values.SpacerSize2)),
			widgets.NewVSpacer(values.SpacerSize2),
		), historyPage.Window.Canvas())

		infoPopUp.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(infoIcon).Add(fyne.NewPos(10, infoIcon.MinSize().Height+5).Subtract(fyne.NewPos(infoPopUp.MinSize().Width, 0))))
	})

	historyPage.HistoryPageContents.Append(widget.NewHBox(historyTitleLabel, layout.NewSpacer(), infoIcon, widgets.NewHSpacer(values.SpacerSize10)))
	return nil
}

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
		infoLabel := widget.NewLabelWithStyle(values.PageHint, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
		gotItLabel := canvas.NewText(values.GotIt, values.Blue)
		gotItLabel.TextStyle = fyne.TextStyle{Bold: true}
		gotItLabel.TextSize = 14

		var infoPopUp *widget.PopUp
		infoPopUp = widget.NewPopUp(widget.NewVBox(
			widgets.NewVSpacer(5),
			widget.NewHBox(widgets.NewHSpacer(5), infoLabel, widgets.NewHSpacer(5)),
			widgets.NewVSpacer(5),
			widget.NewHBox(layout.NewSpacer(), widgets.NewClickableBox(widget.NewVBox(gotItLabel), func() { infoPopUp.Hide() }), widgets.NewHSpacer(5)),
			widgets.NewVSpacer(5),
		), historyPage.Window.Canvas())

		infoPopUp.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(infoIcon).Add(fyne.NewPos(0, infoIcon.Size().Height)))
	})

	historyPage.HistoryPageContents.Append(widget.NewHBox(historyTitleLabel, widgets.NewHSpacer(110), infoIcon))
	return nil
}

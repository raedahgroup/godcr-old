package walletshandler

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (walletPage *WalletPageObject) dialogMenu(posOfIcon fyne.Position) *widget.PopUp {
	clickableText := func(text string, callFunc func()) *widgets.ClickableWidget {
		TextWithPadding := widget.NewHBox(widgets.NewHSpacer(values.SpacerSize12), widgets.NewTextWithSize(text, values.DefaultTextColor, 14), layout.NewSpacer(), widgets.NewHSpacer(values.SpacerSize40))
		textBox := widget.NewVBox(
			widgets.NewVSpacer(values.SpacerSize12),
			TextWithPadding,
			widgets.NewVSpacer(values.SpacerSize12),
		)

		return widgets.NewClickableWidget(textBox, callFunc)
	}
	var popUp *widget.PopUp
	callFunc := func() {
		popUp.Hide()
	}

	dialogBox := widget.NewVBox(
		widgets.NewHSpacer(values.SpacerSize4),
		clickableText(values.SignMessage, callFunc),
		clickableText(values.VerifyMessage, callFunc),
		widgets.NewHSpacer(values.SpacerSize4),
		canvas.NewLine(values.StrippedLineColor),
		widgets.NewHSpacer(values.SpacerSize4),
		clickableText(values.ViewProperty, callFunc),
		widgets.NewHSpacer(values.SpacerSize4),
		canvas.NewLine(values.StrippedLineColor),
		widgets.NewHSpacer(values.SpacerSize4),
		clickableText(values.RenameWallet, callFunc),
		clickableText(values.WalletSettings, callFunc),
		widgets.NewHSpacer(values.SpacerSize4),
	)

	posX := dialogBox.MinSize().Width

	popUp = widget.NewPopUpAtPosition(dialogBox, walletPage.Window.Canvas(), posOfIcon.Subtract(fyne.NewPos(posX, 0).Subtract(fyne.NewPos(0, 20))))
	return popUp
}

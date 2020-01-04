package walletshandler

import (
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (walletPage *WalletPageObject) dialogMenu(walletLabel *canvas.Text, posOfIcon fyne.Position, walletID int) *widget.PopUp {
	var popUp *widget.PopUp

	clickableText := func(text string, callFunc func()) *widgets.ClickableWidget {
		TextWithPadding := widget.NewHBox(widgets.NewHSpacer(values.SpacerSize12), widgets.NewTextWithSize(text, values.DefaultTextColor, 14), layout.NewSpacer(), widgets.NewHSpacer(values.SpacerSize40))
		textBox := widget.NewVBox(
			widgets.NewVSpacer(values.SpacerSize12),
			TextWithPadding,
			widgets.NewVSpacer(values.SpacerSize12),
		)

		return widgets.NewClickableWidget(textBox, callFunc)
	}

	callFunc := func() {
		popUp.Hide()
	}

	renameWalletFunc := func() {
		walletPage.renameWalletPopUp(walletID, walletLabel)
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
		clickableText(values.RenameWallet, renameWalletFunc),
		clickableText(values.WalletSettings, callFunc),
		widgets.NewHSpacer(values.SpacerSize4),
	)

	posX := dialogBox.MinSize().Width

	popUp = widget.NewPopUpAtPosition(dialogBox, walletPage.Window.Canvas(), posOfIcon.Subtract(fyne.NewPos(posX, 0).Subtract(fyne.NewPos(0, 20))))
	return popUp
}

func (walletPage *WalletPageObject) renameWalletPopUp(walletID int, walletLabel *canvas.Text) { //baseText string, onRename func(string) error, onCancel func(*widget.PopUp), otherCallFunc func(string)) {
	onRename := func(value string) error {
		return walletPage.MultiWallet.RenameWallet(walletID, value)
	}
	onCancel := func(popup *widget.PopUp) {
		popup.Hide()
	}
	otherCallFunc := func(value string) {
		walletLabel.Text = value
		walletPage.showLabel("Wallet renamed", walletPage.successLabel)
	}

	walletPage.renameAccountOrWalletPopUp(values.RenameWallet, onRename, onCancel, otherCallFunc)
}

func (walletPage *WalletPageObject) showLabel(Text string, object *widgets.BorderedText) {
	object.SetText(Text)
	object.Container.Show()
	walletPage.WalletPageContents.Refresh()
	time.AfterFunc(time.Second*5, func() {
		object.Container.Hide()
		walletPage.WalletPageContents.Refresh()
	})
}

package pages

import (
	"context"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

//exit functions blocks interactions till user exit goes back to overview if user doesnt
func exit(ctx context.Context, app fyne.App, window fyne.Window) fyne.CanvasObject {
	var popup *widget.PopUp

	yesButton := widget.NewButtonWithIcon("Yes", theme.ConfirmIcon(), func() {
		window.Close()
		<-ctx.Done()
		app.Quit()
	})
	noButton := widget.NewButtonWithIcon("No", theme.CancelIcon(), func() {
		menu.tabs.SelectTabIndex(0)
		popup.Hide()
	})

	exitView := widget.NewVBox(widget.NewLabelWithStyle("Exit", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true}), widget.NewHBox(yesButton, noButton))
	popup = widget.NewModalPopUp(exitView, window.Canvas())
	return popup
}
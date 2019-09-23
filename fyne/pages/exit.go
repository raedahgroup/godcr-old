package pages

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

//exit functions blocks interactions till user exit goes back to overview if user doesnt
func (app *AppInterface) exitPageContent() fyne.CanvasObject {
	var popup *widget.PopUp

	yesButton := widget.NewButtonWithIcon("Yes", theme.ConfirmIcon(), func() {
		app.Window.Close()
		fyne.CurrentApp().Quit()
		fmt.Println("Exited fyne")
	})
	noButton := widget.NewButtonWithIcon("No", theme.CancelIcon(), func() {
		app.tabMenu.SelectTabIndex(0)
		popup.Hide()
	})

	exitView := widget.NewVBox(widget.NewLabelWithStyle("Exit", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true}), widget.NewHBox(yesButton, noButton))
	popup = widget.NewModalPopUp(exitView, app.Window.Canvas())
	return popup
}

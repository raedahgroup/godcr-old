package pages

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (app *AppInterface) exitPageContent() fyne.CanvasObject {
	var popup *widget.PopUp

	yesButton := widget.NewButtonWithIcon("Yes", theme.ConfirmIcon(), func() {
		app.Window.Close()
		fmt.Println("Exited fyne")
	})
	noButton := widget.NewButtonWithIcon("No", theme.CancelIcon(), func() {
		app.tabMenu.SelectTabIndex(0)
		popup.Hide()
	})

	exitView := widget.NewVBox(
		widgets.NewVSpacer(values.Padding),
		widget.NewLabelWithStyle("Exit?", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widgets.NewVSpacer(values.Padding),
		widget.NewHBox(yesButton, widgets.NewHSpacer(values.SpacerSize4), noButton),
		widgets.NewVSpacer(values.Padding))

	viewWithPadding := widget.NewHBox(
		widgets.NewHSpacer(values.Padding), exitView, widgets.NewHSpacer(values.Padding))

	popup = widget.NewModalPopUp(viewWithPadding, app.Window.Canvas())
	popup.Show()
	return popup
}

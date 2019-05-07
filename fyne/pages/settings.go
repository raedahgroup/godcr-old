package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

//SettingsPage comprises of all setting functions they are to be created and passed to the menu function
func SettingsPage(win fyne.Window, App fyne.App) fyne.CanvasObject {
	theme := changeTheme(win, App)
	return Menu(theme, win, App)
}

func changeTheme(win fyne.Window, change fyne.App) fyne.CanvasObject {
	var radio1 *widget.Radio
	var radio2 *widget.Radio
	radio1 = widget.NewRadio([]string{"Light Theme"}, func(background string) {
		if background == "Light Theme" {
			if radio2.Selected == "Dark Theme" {
				radio2.SetSelected("")
			}
			change.Settings().SetTheme(theme.LightTheme())
		}
	})
	radio2 = widget.NewRadio([]string{"Dark Theme"}, func(background string) {
		if background == "Dark Theme" {
			if radio1.Selected == "Light Theme" {
				radio1.SetSelected("")
			}
			change.Settings().SetTheme(theme.DarkTheme())
		}
	})
	orderedRadio := []fyne.CanvasObject{
		radio1, radio2,
	}
	radio := layout.NewHBoxLayout()
	radio.Layout(orderedRadio, fyne.NewSize(0, 0))
	return fyne.NewContainerWithLayout(radio, radio1, radio2)
}

func changeSpendingPass(App fyne.App) {

	// PassphraseWindow := App.NewWindow("Enter Current Pin")
	// widget.NewForm()
}

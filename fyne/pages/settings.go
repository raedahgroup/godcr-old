package pages

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func settingsPage(App fyne.App) fyne.CanvasObject {
	settings := widget.NewVBox(changeTheme(App))
	return widget.NewHBox(widgets.NewHSpacer(10), settings)
}

//todo: after changing theme, make it default
func changeTheme(change fyne.App) fyne.CanvasObject {
	fyneTheme := change.Settings().Theme()

	radio := widget.NewRadio([]string{"Light Theme", "Dark Theme"}, func(background string) {
		if background == "Light Theme" {
			change.Settings().SetTheme(theme.LightTheme())
			receive.errorLabel.Color = color.RGBA{255, 0, 0, 255}
			menu.alphaTheme = 255
		} else if background == "Dark Theme" {
			change.Settings().SetTheme(theme.DarkTheme())
			receive.errorLabel.Color = color.RGBA{255, 0, 0, 0}
			menu.alphaTheme = 200
		}
	})
	radio.Horizontal = true

	if fyneTheme.BackgroundColor() == theme.LightTheme().BackgroundColor() {
		radio.SetSelected("Light Theme")
	} else if fyneTheme.BackgroundColor() == theme.DarkTheme().BackgroundColor() {
		radio.SetSelected("Dark Theme")
	}
	return radio
}

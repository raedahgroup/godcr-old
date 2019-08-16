package pages

import (
	"image/color"
	"io/ioutil"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func settingsPage(App fyne.App) fyne.CanvasObject {
	settings := widget.NewVBox(changeTheme(App))
	return widget.NewHBox(widgets.NewHSpacer(10), settings)
}

// todo: after changing theme, make it default
func changeTheme(change fyne.App) fyne.CanvasObject {
	fyneTheme := change.Settings().Theme()

	radio := widget.NewRadio([]string{"Light Theme", "Dark Theme"}, func(background string) {
		if background == "Light Theme" {
			// set overview icon and name
			decredDark, err := ioutil.ReadFile("./fyne/pages/png/decredDark.png")
			if err != nil {
				log.Fatalln("exit png file missing", err)
			}
			overview.goDcrLabel.Color = color.RGBA{0, 0, 0, 255}
			iconResource := canvas.NewImageFromResource(fyne.NewStaticResource("Decred", decredDark))
			overview.icon.Resource = iconResource.Resource
			canvas.Refresh(overview.icon)

			change.Settings().SetTheme(theme.LightTheme())
		} else if background == "Dark Theme" {
			decredLight, err := ioutil.ReadFile("./fyne/pages/png/decredLight.png")
			if err != nil {
				log.Fatalln("exit png file missing", err)
			}

			overview.goDcrLabel.Color = color.RGBA{255, 255, 255, 0}
			iconResource := canvas.NewImageFromResource(fyne.NewStaticResource("Decred", decredLight))
			overview.icon.Resource = iconResource.Resource
			canvas.Refresh(overview.icon)
			canvas.Refresh(overview.goDcrLabel)

			change.Settings().SetTheme(theme.DarkTheme())
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

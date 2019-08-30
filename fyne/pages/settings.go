package pages

import (
	"image/color"
	"io/ioutil"
	"log"

	"github.com/raedahgroup/dcrlibwallet"
	godcrApp "github.com/raedahgroup/godcr/app"
	dlw "github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/widgets"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func settingsPage(App fyne.App, wallet godcrApp.WalletMiddleware) fyne.CanvasObject {
	settings := widget.NewVBox(
		widgets.NewItalicizedLabel("Staking"),
		stakingForm(wallet),
		widgets.NewVSpacer(10),
		changeTheme(App),
	)

	return widget.NewHBox(
		widgets.NewHSpacer(10),
		settings,
	)
}

func stakingForm(wallet godcrApp.WalletMiddleware) fyne.Widget {
	libwallet, ok := wallet.(*dlw.DcrWalletLib)
	if !ok {
		return widgets.NewItalicizedLabel("Settings not yet supported for dcrwallet")
	}

	var vspHostValue string
	if err := libwallet.ReadFromSettings(dcrlibwallet.VSPHostSettingsKey, &vspHostValue); err != nil {
		println(err.Error())
	}

	vspHost := widget.NewEntry()
	vspHost.SetText(vspHostValue)

	stakingForm := &widget.Form{
		OnSubmit: func() {
			if err := libwallet.SaveToSettings(dcrlibwallet.VSPHostSettingsKey, vspHost.Text); err != nil {
				println(err.Error())
			}
		},
	}
	stakingForm.Append("VSP Host", vspHost)

	return stakingForm
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

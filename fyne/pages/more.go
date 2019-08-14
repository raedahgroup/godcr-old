package pages

import (
	"io/ioutil"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func morePage(fyneApp fyne.App) fyne.CanvasObject {
	settingFile, err := ioutil.ReadFile("./fyne/pages/png/settings.png")
	if err != nil {
		log.Fatalln("setting png file missing", err)
	}
	helpFile, err := ioutil.ReadFile("./fyne/pages/png/help.png")
	if err != nil {
		log.Fatalln("help png file missing", err)
	}
	securityFile, err := ioutil.ReadFile("./fyne/pages/png/security.png")
	if err != nil {
		log.Fatalln("security png file missing", err)
	}
	infoFile, err := ioutil.ReadFile("./fyne/pages/png/info.png")
	if err != nil {
		log.Fatalln("info png file missing", err)
	}
	stakingFile, err := ioutil.ReadFile("./fyne/pages/png/stake.png")
	if err != nil {
		log.Fatalln("staking png file missing", err)
	}

	container := widget.NewTabContainer(

		widget.NewTabItemWithIcon("Staking", fyne.NewStaticResource("Staking", stakingFile), widget.NewLabel("Hello")),
		widget.NewTabItemWithIcon("Settings", fyne.NewStaticResource("More", settingFile), settingsPage(fyneApp)),
		widget.NewTabItemWithIcon("Security Tools", fyne.NewStaticResource("More", securityFile), widget.NewLabel("Hello")),
		widget.NewTabItemWithIcon("Help", fyne.NewStaticResource("More", helpFile), widget.NewLabel("Hello")),
		widget.NewTabItemWithIcon("Debug", fyne.NewStaticResource("More", settingFile), widget.NewLabel("Hello")),
		widget.NewTabItemWithIcon("About", fyne.NewStaticResource("More", infoFile), widget.NewLabel("Hello")))
	return container
}

package pages

import (
	"io/ioutil"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
)

func morePage(wallet godcrApp.WalletMiddleware, fyneApp fyne.App) fyne.CanvasObject {
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
	aboutFile, err := ioutil.ReadFile("./fyne/pages/png/about.png")
	if err != nil {
		log.Fatalln("about png file missing", err)
	}

	container := widget.NewTabContainer(
		widget.NewTabItemWithIcon("Settings", fyne.NewStaticResource("More", settingFile), settingsPage(fyneApp)),
		widget.NewTabItemWithIcon("Security Tools", fyne.NewStaticResource("More", securityFile), pageNotImplemented()),
		widget.NewTabItemWithIcon("Help", fyne.NewStaticResource("More", helpFile), pageNotImplemented()),
		widget.NewTabItemWithIcon("Debug", fyne.NewStaticResource("More", settingFile), pageNotImplemented()),
		widget.NewTabItemWithIcon("About", fyne.NewStaticResource("More", aboutFile), pageNotImplemented()))
	return container
}

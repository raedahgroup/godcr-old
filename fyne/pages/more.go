package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
)

func morePage(wallet godcrApp.WalletMiddleware, fyneApp fyne.App) fyne.CanvasObject {
	return widget.NewTabContainer(
		widget.NewTabItemWithIcon("Settings", settingsIcon, settingsPage(fyneApp)),
		widget.NewTabItemWithIcon("Security Tools", securityIcon, pageNotImplemented()),
		widget.NewTabItemWithIcon("Help", helpIcon, pageNotImplemented()),
		widget.NewTabItemWithIcon("Debug", settingsIcon, pageNotImplemented()),
		widget.NewTabItemWithIcon("About", aboutIcon, pageNotImplemented()))
}

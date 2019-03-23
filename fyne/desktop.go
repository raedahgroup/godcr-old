package fyne

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func LaunchApp() {
	a := app.New()
	w := a.NewWindow("GoDCR Wallet")

	content := widget.NewVBox(
	)

	menu := widget.NewVBox(
		widget.NewGroup("GoDCR", widget.NewVBox(
			widget.NewButton("Dashboard", func() {

			}),
			widget.NewButton("Accounts", func() {

			}),
			widget.NewButton("Send", func() {

			}),
			widget.NewButton("Recieve", func() {

			}),
			widget.NewButton("Quit", func() {
				a.Quit()
			}),
		)),
	)

	w.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			fyne.NewContainerWithLayout(layout.NewGridLayout(1),
				menu,
			),
			fyne.NewContainerWithLayout(layout.NewGridLayout(8),
				content,
			),
		),
	))

	w.ShowAndRun()
}

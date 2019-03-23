package fyne

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func LaunchApp() {
	a := app.New()

	w := a.NewWindow("Hello")
	w.FixedSize()
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			a.Quit()
		}),
	))

	w.ShowAndRun()
}

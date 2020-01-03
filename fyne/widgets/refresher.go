package widgets

import (
	"fyne.io/fyne"
)

func Refresher(objects ...fyne.CanvasObject) {
	for _, object := range objects {
		object.Refresh()
	}
}

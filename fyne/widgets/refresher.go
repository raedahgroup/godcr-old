package widgets

import (
	"fyne.io/fyne"
)

func Refresher(objects ...fyne.Widget) {
	for _, object := range objects {
		object.Refresh()
	}
}

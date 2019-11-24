package layouts

import (
	"fyne.io/fyne"
)

type boxLayout struct {
	horizontal bool
	spacer     int
}

func (c *boxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 1 {
		return
	}
	objects[0].Move(fyne.NewPos(0, 0))
	objects[0].Resize(objects[0].MinSize())

	for index, object := range objects {
		if index == 0 {
			continue
		}

		if c.horizontal {
			redH := objects[index-1].MinSize().Height - object.MinSize().Height
			if redH < 0 {
				redH = 0
			}
			object.Move(fyne.NewPos(objects[index-1].Position().X+objects[index-1].MinSize().Width+c.spacer, objects[index-1].Position().Y+redH-1)) //, objects[index-1].Position().Y+object.MinSize().Height)))
			object.Resize(object.MinSize())
		} else {
			object.Move(fyne.NewPos(objects[index-1].Position().X+5, objects[index-1].Position().Y+object.MinSize().Height+c.spacer-10))
			object.Resize(object.MinSize())
		}
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
func (c *boxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var minSize fyne.Size

	for _, child := range objects {
		if c.horizontal {
			minSize = minSize.Add(fyne.NewSize(child.MinSize().Width+c.spacer, 0))
			minSize.Height = fyne.Max(child.MinSize().Height, minSize.Height)
		} else {
			minSize = minSize.Add(fyne.NewSize(0-5, child.MinSize().Height+c.spacer-6))
			minSize.Width = fyne.Max(child.MinSize().Width, minSize.Width)
		}
	}

	return minSize
}

func NewHBox(spacer int) fyne.Layout {
	return &boxLayout{true, spacer}
}

func NewVBox(spacer int) fyne.Layout {
	return &boxLayout{false, spacer}
}

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
	objects[0].Resize(objects[0].MinSize())
	objects[0].Move(fyne.NewPos(0, 0))

	for index, object := range objects {
		if index == 0 {
			continue
		}

		if c.horizontal {
			object.Resize(object.MinSize())
			object.Move(fyne.NewPos(objects[index-1].Position().X+object.MinSize().Width+c.spacer, objects[index-1].Position().Y+5))
		} else {
			object.Resize(object.MinSize())
			object.Move(fyne.NewPos(objects[index-1].Position().X+5, objects[index-1].Position().Y+object.MinSize().Height+c.spacer))
		}
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
func (c *boxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var minSize fyne.Size

	if c.horizontal {
		minSize = minSize.Add(fyne.NewSize(6-c.spacer, 0))
	} else {
		minSize = minSize.Add(fyne.NewSize(0, 6-c.spacer))
	}

	for _, child := range objects {
		if c.horizontal {
			minSize = minSize.Add(fyne.NewSize(child.MinSize().Width+c.spacer-6, 0-5))
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

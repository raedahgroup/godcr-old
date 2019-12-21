package layouts

import (
	"fyne.io/fyne"
)

type boxLayout struct {
	horizontal bool
	spacer     int
	isAmount   bool
}

func (c *boxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 1 {
		return
	}

	for index, object := range objects {
		if index == 0 {
			continue
		}

		if c.horizontal {
			hSize := objects[index-1].MinSize().Height - object.MinSize().Height
			if c.isAmount {
				hSize = hSize - 3
			}
			object.Move(fyne.NewPos(objects[index-1].Position().X+objects[index-1].MinSize().Width+c.spacer, hSize))
		} else {
			hSize := objects[index-1].MinSize().Width - object.MinSize().Width
			if c.isAmount {
				hSize = hSize - 3
			}
			object.Move(fyne.NewPos(hSize, objects[index-1].Position().Y+objects[index-1].MinSize().Height+c.spacer))
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

func NewHBox(spacer int, isAmount bool) fyne.Layout {
	return &boxLayout{true, spacer, isAmount}
}

func NewVBox(spacer int) fyne.Layout {
	return &boxLayout{false, spacer, false}
}

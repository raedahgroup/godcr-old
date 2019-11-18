package layouts

import "fyne.io/fyne"

type passwordLayout struct {
	passwordSize fyne.Size
}

// Layout places the icon in the password entry. password should be the first index
func (c *passwordLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	password := objects[0]
	icon := objects[1]

	password.Move(fyne.NewPos(0, 0))
	password.Resize(c.passwordSize)

	icon.Move(fyne.NewPos(password.Position().X+c.passwordSize.Width-icon.MinSize().Width, password.Position().Y))
	icon.Resize(icon.MinSize())
}

// MinSize finds the smallest size that satisfies all the child objects.
func (c *passwordLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return c.passwordSize
}

func NewPasswordLayout(passwordSize fyne.Size) fyne.Layout {
	return &passwordLayout{passwordSize}
}

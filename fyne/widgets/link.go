package widgets

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"github.com/raedahgroup/godcr/fyne/styles"
)

var (
	defaultTheme = styles.NewTheme()
)

type linkRenderer struct {
	icon  *canvas.Image
	label *canvas.Text

	objects []fyne.CanvasObject
	link    *Link
}

// MinSize calculates the minimum size of a button.
// This is based on the contained text, any icon that is set and a standard
// amount of padding added.
func (l *linkRenderer) MinSize() fyne.Size {
	return l.label.MinSize().Add(fyne.NewSize(defaultTheme.Padding()*14, defaultTheme.Padding()*2))
}

// Layout the components of the button widget
func (l *linkRenderer) Layout(size fyne.Size) {
	inner := size.Subtract(fyne.NewSize(defaultTheme.Padding()*4, defaultTheme.Padding()*2))
	l.label.Resize(inner)
	l.label.Move(fyne.NewPos(defaultTheme.Padding()*2, defaultTheme.Padding()))

}

// ApplyTheme is called when the Button may need to update it's look
func (l *linkRenderer) ApplyTheme() {
	l.label.Color = defaultTheme.TextColor()

	l.Refresh()
}

func (l *linkRenderer) BackgroundColor() color.Color {
	return defaultTheme.BackgroundColor()
}

func (l *linkRenderer) Refresh() {
	l.label.Text = l.link.Text
	l.Layout(l.link.Size())
	canvas.Refresh(l.link)
}

func (l *linkRenderer) Objects() []fyne.CanvasObject {
	return l.objects
}

func (l *linkRenderer) Destroy() {
}

// Button widget has a text label and triggers an event func when clicked
type Link struct {
	baseWidget
	Text  string
	Style TextStyle

	OnTapped func() `json:"-"`
}

// ButtonStyle determines the behaviour and rendering of a button.
type TextStyle int

const (
	// DefaultButton is the standard button style
	DefaultButton TextStyle = iota
	// PrimaryButton that should be more prominent to the user
	PrimaryButton
)

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (l *Link) Resize(size fyne.Size) {
	l.resize(size, l)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (l *Link) Move(pos fyne.Position) {
	l.move(pos, l)
}

// MinSize returns the smallest size this widget can shrink to
func (l *Link) MinSize() fyne.Size {
	return l.minSize(l)
}

// Show this widget, if it was previously hidden
func (l *Link) Show() {
	l.show(l)
}

// Hide this widget, if it was previously visible
func (l *Link) Hide() {
	l.hide(l)
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (l *Link) Tapped(*fyne.PointEvent) {
	if l.OnTapped != nil {
		l.OnTapped()
	}
}

// TappedSecondary is called when a secondary pointer tapped event is captured
func (l *Link) TappedSecondary(*fyne.PointEvent) {
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (l *Link) CreateRenderer() fyne.WidgetRenderer {
	var icon *canvas.Image

	text := canvas.NewText(l.Text, defaultTheme.HyperlinkColor())
	text.Alignment = fyne.TextAlignLeading

	objects := []fyne.CanvasObject{
		text,
	}
	if icon != nil {
		objects = append(objects, icon)
	}

	return &linkRenderer{icon, text, objects, l}
}

// SetText allows the button label to be changed
func (l *Link) SetText(text string) {
	l.Text = text

	Refresh(l)
}

// NewButton creates a new button widget with the set label and tap handler
func NewLink(label string, tapped func()) *Link {
	link := &Link{baseWidget{}, label, DefaultButton, tapped}

	Renderer(link).Layout(link.MinSize())
	return link
}

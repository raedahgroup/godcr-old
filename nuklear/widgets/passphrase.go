package widgets

import (
	"image"
	"image/color"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/rect"
)

type Passphrase struct {
	input        nucular.TextEditor
	bounds       rect.Rect
	headerBounds rect.Rect
	errStr       string
	passphrase   chan string
}

func NewPassphraseWidget() *Passphrase {
	bounds := rect.Rect{
		X: 250,
		Y: 60,
		W: 260,
		H: 150,
	}

	headerBounds := rect.Rect{
		X: 10,
		Y: 10,
		W: bounds.W,
		H: 30,
	}

	passphraseWidget := &Passphrase{
		bounds:       bounds,
		headerBounds: headerBounds,
	}
	passphraseWidget.input.Flags = nucular.EditSimple
	passphraseWidget.input.PasswordChar = '*'

	return passphraseWidget
}

func (p *Passphrase) Get(window *nucular.Window, passphrase chan string) {
	p.passphrase = passphrase
	window.Master().PopupOpen("Wallet Passphrase", nucular.WindowTitle|nucular.WindowDynamic|nucular.WindowNoScrollbar, p.bounds, true, p.popup)
}

func (p *Passphrase) popup(window *nucular.Window) {
	masterWindow := window.Master()

	// set popup style
	style := window.Master().Style()
	style.NormalWindow.Padding = image.Point{20, 50}
	style.NormalWindow.Background = color.RGBA{0xff, 0xff, 0xff, 0xff}
	masterWindow.SetStyle(style)

	defer func() {
		// reset page style
		style.NormalWindow.Padding = image.Point{0, 0}
		masterWindow.SetStyle(style)
	}()

	// render popup
	window.Row(20).Dynamic(1)
	window.Label("Passphrase", "LC")

	window.Row(25).Dynamic(1)
	p.input.Edit(window)

	if p.errStr != "" {
		window.Row(10).Dynamic(1)
		window.LabelColored(p.errStr, "LC", color.RGBA{205, 32, 32, 255})
	}

	window.Row(5).Dynamic(1)
	window.Label("", "LC")

	window.Row(25).Static(65, 65)
	if window.Button(label.T("Close"), false) {
		window.Close()
	}

	if window.Button(label.T("Submit"), false) {
		if p.validate() {
			p.passphrase <- string(p.input.Buffer)
			window.Close()
			return
		}
		masterWindow.Changed()
	}
}

func (p *Passphrase) validate() bool {
	if string(p.input.Buffer) == "" {
		p.errStr = "Wallet passphrase is required"
		return false
	}
	p.errStr = ""
	return true
}

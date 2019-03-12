package widgets

import (
	"reflect"
	"testing"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
)

func TestNewPassphraseWidget(t *testing.T) {
	tests := []struct {
		name string
		want *Passphrase
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := NewPassphraseWidget(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("NewPassphraseWidget() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestPassphrase_Get(t *testing.T) {
	type fields struct {
		input        nucular.TextEditor
		bounds       rect.Rect
		headerBounds rect.Rect
		errStr       string
		passphrase   chan string
	}
	tests := []struct {
		name       string
		fields     fields
		window     *nucular.Window
		passphrase chan string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &Passphrase{
				input:        test.fields.input,
				bounds:       test.fields.bounds,
				headerBounds: test.fields.headerBounds,
				errStr:       test.fields.errStr,
				passphrase:   test.fields.passphrase,
			}
			p.Get(test.window, test.passphrase)
		})
	}
}

func TestPassphrase_popup(t *testing.T) {
	type fields struct {
		input        nucular.TextEditor
		bounds       rect.Rect
		headerBounds rect.Rect
		errStr       string
		passphrase   chan string
	}
	tests := []struct {
		name   string
		fields fields
		window *nucular.Window
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &Passphrase{
				input:        test.fields.input,
				bounds:       test.fields.bounds,
				headerBounds: test.fields.headerBounds,
				errStr:       test.fields.errStr,
				passphrase:   test.fields.passphrase,
			}
			p.popup(test.window)
		})
	}
}

func TestPassphrase_validate(t *testing.T) {
	type fields struct {
		input        nucular.TextEditor
		bounds       rect.Rect
		headerBounds rect.Rect
		errStr       string
		passphrase   chan string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &Passphrase{
				input:        test.fields.input,
				bounds:       test.fields.bounds,
				headerBounds: test.fields.headerBounds,
				errStr:       test.fields.errStr,
				passphrase:   test.fields.passphrase,
			}
			if got := p.validate(); got != test.want {
				t.Errorf("Passphrase.validate() = %v, want %v", got, test.want)
			}
		})
	}
}

package helpers

import (
	"reflect"
	"testing"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
)

func TestNewWindow(t *testing.T) {
	tests := []struct {
		name  string
		title string
		w     *nucular.Window
		flags nucular.WindowFlags
		want  *Window
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := NewWindow(test.title, test.w, test.flags); !reflect.DeepEqual(got, test.want) {
				t.Errorf("NewWindow() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWindow_DrawHeader(t *testing.T) {
	type fields struct {
		Window *nucular.Window
	}
	tests := []struct {
		name   string
		fields fields
		title  string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := &Window{
				Window: test.fields.Window,
			}
			w.DrawHeader(test.title)
		})
	}
}

func TestWindow_ContentWindow(t *testing.T) {
	type fields struct {
		Window *nucular.Window
	}
	tests := []struct {
		name   string
		fields fields
		title  string
		want   *Window
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := &Window{
				Window: test.fields.Window,
			}
			if got := w.ContentWindow(test.title); !reflect.DeepEqual(got, test.want) {
				t.Errorf("Window.ContentWindow() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWindow_SetErrorMessage(t *testing.T) {
	type fields struct {
		Window *nucular.Window
	}
	tests := []struct {
		name    string
		fields  fields
		message string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := &Window{
				Window: test.fields.Window,
			}
			w.SetErrorMessage(test.message)
		})
	}
}

func TestWindow_Style(t *testing.T) {
	type fields struct {
		Window *nucular.Window
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := &Window{
				Window: test.fields.Window,
			}
			w.Style()
		})
	}
}

func TestWindow_End(t *testing.T) {
	type fields struct {
		Window *nucular.Window
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := &Window{
				Window: test.fields.Window,
			}
			w.End()
		})
	}
}

func TestSetContentArea(t *testing.T) {
	tests := []struct {
		name string
		area rect.Rect
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			SetContentArea(test.area)
		})
	}
}

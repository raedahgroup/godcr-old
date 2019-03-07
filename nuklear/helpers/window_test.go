package helpers

import (
	"reflect"
	"testing"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
)

func TestNewWindow(t *testing.T) {
	tests := []struct {
		name   string
		title  string
		window *nucular.Window
		flags  nucular.WindowFlags
		want   *Window
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := NewWindow(test.title, test.window, test.flags); !reflect.DeepEqual(got, test.want) {
				t.Errorf("NewWindow() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWindow_DrawHeader(t *testing.T) {
	tests := []struct {
		name   string
		window *Window
		title  string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.window.DrawHeader(test.title)
		})
	}
}

func TestWindow_ContentWindow(t *testing.T) {
	tests := []struct {
		name   string
		window *Window
		title  string
		want   *Window
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.window.ContentWindow(test.title); !reflect.DeepEqual(got, test.want) {
				t.Errorf("Window.ContentWindow() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestWindow_SetErrorMessage(t *testing.T) {
	tests := []struct {
		name    string
		window  *Window
		message string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.window.SetErrorMessage(test.message)
		})
	}
}

func TestWindow_Style(t *testing.T) {
	tests := []struct {
		name   string
		window *Window
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.window.Style()
		})
	}
}

func TestWindow_End(t *testing.T) {
	tests := []struct {
		name   string
		window *Window
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.window.End()
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

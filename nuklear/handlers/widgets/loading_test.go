package widgets

import (
	"reflect"
	"testing"

	"github.com/aarzilli/nucular"
)

func TestNewLoadingWidget(t *testing.T) {
	tests := []struct {
		name string
		want Widget
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := NewLoadingWidget(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("NewLoadingWidget() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestLoadingWidget_Render(t *testing.T) {
	tests := []struct {
		name   string
		l      *LoadingWidget
		window *nucular.Window
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &LoadingWidget{}
			l.Render(test.window)
		})
	}
}

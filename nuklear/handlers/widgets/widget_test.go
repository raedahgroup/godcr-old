package widgets

import (
	"testing"

	"github.com/aarzilli/nucular"
)

func TestShowLoadingWidget(t *testing.T) {
	tests := []struct {
		name   string
		window *nucular.Window
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ShowLoadingWidget(test.window)
		})
	}
}

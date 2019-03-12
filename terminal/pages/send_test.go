package pages

import (
	"reflect"
	"testing"

	"github.com/rivo/tview"
)

func TestSendPage(t *testing.T) {
	tests := []struct {
		name       string
		setFocus   func(p tview.Primitive) *tview.Application
		clearFocus func()
		want       tview.Primitive
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := SendPage(test.setFocus, test.clearFocus); !reflect.DeepEqual(got, test.want) {
				t.Errorf("SendPage() = %v, want %v", got, test.want)
			}
		})
	}
}

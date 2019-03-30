package widgets

import (
	"reflect"
	"testing"
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

package nuklear

import (
	"reflect"
	"testing"
)

func Test_getNavPages(t *testing.T) {
	tests := []struct {
		name string
		want []navPage
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := getNavPages(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("getNavPages() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_getStandalonePages(t *testing.T) {
	tests := []struct {
		name string
		want []standalonePage
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := getStandalonePages(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("getStandalonePages() = %v, want %v", got, test.want)
			}
		})
	}
}

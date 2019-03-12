package routes

import (
	"html/template"
	"reflect"
	"testing"
)

func Test_templates(t *testing.T) {
	tests := []struct {
		name string
		want []templateData
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := templates(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("templates() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_templateFuncMap(t *testing.T) {
	tests := []struct {
		name string
		want template.FuncMap
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := templateFuncMap(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("templateFuncMap() = %v, want %v", got, test.want)
			}
		})
	}
}

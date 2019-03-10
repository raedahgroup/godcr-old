package commands

import (
	"reflect"
	"testing"

	"github.com/raedahgroup/godcr/app/help"
)

func TestCategories(t *testing.T) {
	tests := []struct {
		name string
		want []*help.CommandCategory
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Categories(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("Categories() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_commanderStub_Execute(t *testing.T) {
	tests := []struct {
		name    string
		w       commanderStub
		args    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := commanderStub{}
			if err := w.Execute(test.args); (err != nil) != test.wantErr {
				t.Errorf("commanderStub.Execute() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

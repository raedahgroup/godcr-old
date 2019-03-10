package commands

import (
	"reflect"
	"testing"

	"github.com/jessevdk/go-flags"
)

func TestHelpCommand_Run(t *testing.T) {
	type fields struct {
		commanderStub commanderStub
		Args          struct{ CommandName string }
	}
	tests := []struct {
		name    string
		fields  fields
		parser  *flags.Parser
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := HelpCommand{
				commanderStub: test.fields.commanderStub,
				Args:          test.fields.Args,
			}
			if err := h.Run(test.parser); (err != nil) != test.wantErr {
				t.Errorf("HelpCommand.Run() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestHelpParser(t *testing.T) {
	tests := []struct {
		name string
		want *flags.Parser
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := HelpParser(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("HelpParser() = %v, want %v", got, test.want)
			}
		})
	}
}

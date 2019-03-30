package help

import (
	"bytes"
	"testing"

	flags "github.com/jessevdk/go-flags"
)

func Test_commandCategoryName(t *testing.T) {
	tests := []struct {
		name              string
		commandName       string
		commandCategories []*CommandCategory
		want              string
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := commandCategoryName(test.commandName, test.commandCategories); got != test.want {
				t.Errorf("commandCategoryName() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestPrintGeneralHelp(t *testing.T) {
	tests := []struct {
		name              string
		parser            *flags.Parser
		commandCategories []*CommandCategory
		wantOutput        string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			PrintGeneralHelp(output, test.parser, test.commandCategories)
			if gotOutput := output.String(); gotOutput != test.wantOutput {
				t.Errorf("PrintGeneralHelp() = %v, want %v", gotOutput, test.wantOutput)
			}
		})
	}
}

func TestPrintCommandHelp(t *testing.T) {
	tests := []struct {
		name       string
		appName    string
		command    *flags.Command
		wantOutput string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			PrintCommandHelp(output, test.appName, test.command)
			if gotOutput := output.String(); gotOutput != test.wantOutput {
				t.Errorf("PrintCommandHelp() = %v, want %v", gotOutput, test.wantOutput)
			}
		})
	}
}

func TestPrintOptionsSimple(t *testing.T) {
	tests := []struct {
		name       string
		groups     []*flags.Group
		wantOutput string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			PrintOptionsSimple(output, test.groups)
			if gotOutput := output.String(); gotOutput != test.wantOutput {
				t.Errorf("PrintOptionsSimple() = %v, want %v", gotOutput, test.wantOutput)
			}
		})
	}
}

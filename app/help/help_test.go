package help

import (
	"bytes"
	"testing"

	flags "github.com/btcsuite/go-flags"
)

func Test_commandCategoryName(t *testing.T) {
	tests := []struct {
		name              string
		commandName       string
		commandCategories []*CommandCategory
		want              string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commandCategoryName(tt.commandName, tt.commandCategories); got != tt.want {
				t.Errorf("commandCategoryName() = %v, want %v", got, tt.want)
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			PrintGeneralHelp(output, tt.parser, tt.commandCategories)
			if gotOutput := output.String(); gotOutput != tt.wantOutput {
				t.Errorf("PrintGeneralHelp() = %v, want %v", gotOutput, tt.wantOutput)
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			PrintCommandHelp(output, tt.appName, tt.command)
			if gotOutput := output.String(); gotOutput != tt.wantOutput {
				t.Errorf("PrintCommandHelp() = %v, want %v", gotOutput, tt.wantOutput)
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			PrintOptionsSimple(output, tt.groups)
			if gotOutput := output.String(); gotOutput != tt.wantOutput {
				t.Errorf("PrintOptionsSimple() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

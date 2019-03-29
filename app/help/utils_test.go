package help

import (
	"bytes"
	"testing"

	flags "github.com/jessevdk/go-flags"
)

func Test_printOptionGroups(t *testing.T) {
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
			printOptionGroups(output, tt.groups)
			if gotOutput := output.String(); gotOutput != tt.wantOutput {
				t.Errorf("printOptionGroups() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func Test_printOptions(t *testing.T) {
	tests := []struct {
		name              string
		optionDescription string
		options           []*flags.Option
		wantTabWriter     string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabWriter := &bytes.Buffer{}
			printOptions(tabWriter, tt.optionDescription, tt.options)
			if gotTabWriter := tabWriter.String(); gotTabWriter != tt.wantTabWriter {
				t.Errorf("printOptions() = %v, want %v", gotTabWriter, tt.wantTabWriter)
			}
		})
	}
}

func Test_parseOptionUsageText(t *testing.T) {
	tests := []struct {
		name                    string
		option                  *flags.Option
		hasOptionsWithShortName bool
		wantOptionUsage         string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOptionUsage := parseOptionUsageText(tt.option, tt.hasOptionsWithShortName); gotOptionUsage != tt.wantOptionUsage {
				t.Errorf("parseOptionUsageText() = %v, want %v", gotOptionUsage, tt.wantOptionUsage)
			}
		})
	}
}

func Test_parseOptionDescription(t *testing.T) {
	tests := []struct {
		name            string
		option          *flags.Option
		wantDescription string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDescription := parseOptionDescription(tt.option); gotDescription != tt.wantDescription {
				t.Errorf("parseOptionDescription() = %v, want %v", gotDescription, tt.wantDescription)
			}
		})
	}
}

func Test_printCommands(t *testing.T) {
	tests := []struct {
		name          string
		commandGroups map[string][]*flags.Command
		wantTabWriter string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabWriter := &bytes.Buffer{}
			printCommands(tabWriter, tt.commandGroups)
			if gotTabWriter := tabWriter.String(); gotTabWriter != tt.wantTabWriter {
				t.Errorf("printCommands() = %v, want %v", gotTabWriter, tt.wantTabWriter)
			}
		})
	}
}

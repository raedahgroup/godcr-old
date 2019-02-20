package config

import (
	"testing"

	flags "github.com/jessevdk/go-flags"
)

func Test_createConfigFile(t *testing.T) {
	tests := []struct {
		name           string
		wantSuccessful bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSuccessful := createConfigFile(); gotSuccessful != tt.wantSuccessful {
				t.Errorf("createConfigFile() = %v, want %v", gotSuccessful, tt.wantSuccessful)
			}
		})
	}
}

func Test_parseConfigFile(t *testing.T) {
	tests := []struct {
		name    string
		parser  *flags.Parser
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseConfigFile(tt.parser); (err != nil) != tt.wantErr {
				t.Errorf("parseConfigFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateConfigFile(t *testing.T) {
	tests := []struct {
		name          string
		option        string
		newValue      interface{}
		removeComment bool
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateConfigFile(tt.option, tt.newValue, tt.removeComment); (err != nil) != tt.wantErr {
				t.Errorf("UpdateConfigFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

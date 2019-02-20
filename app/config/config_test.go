package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	_, _, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load godcr config: %s", err)
	}
}

func Test_defaultFileOptions(t *testing.T) {
	tests := []struct {
		name string
		want ConfFileOptions
	}{
		// TODO: add test cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultFileOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultFileOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultConfig(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasConfigFileOption(t *testing.T) {
	tests := []struct {
		name        string
		unknownArgs []string
		want        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasConfigFileOption(tt.unknownArgs); got != tt.want {
				t.Errorf("hasConfigFileOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configFileOptions(t *testing.T) {
	tests := []struct {
		name        string
		wantOptions []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOptions := configFileOptions(); !reflect.DeepEqual(gotOptions, tt.wantOptions) {
				t.Errorf("configFileOptions() = %v, want %v", gotOptions, tt.wantOptions)
			}
		})
	}
}

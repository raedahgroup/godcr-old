package config

import (
	"reflect"
	"testing"
)

func Test_defaultConfig(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := defaultConfig(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("defaultConfig() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_defaultCommandLineOptions(t *testing.T) {
	tests := []struct {
		name string
		want CommandLineOptions
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := defaultCommandLineOptions(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("defaultCommandLineOptions() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		want1   []string
		wantErr bool
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, got1, err := LoadConfig()
			if (err != nil) != test.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("LoadConfig() got = %v, want %v", got, test.want)
			}
			if !reflect.DeepEqual(got1, test.want1) {
				t.Errorf("LoadConfig() got1 = %v, want %v", got1, test.want1)
			}
		})
	}
}

func Test_hasConfigFileOption(t *testing.T) {
	tests := []struct {
		name            string
		commandLineArgs []string
		want            bool
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := hasConfigFileOption(test.commandLineArgs); got != test.want {
				t.Errorf("hasConfigFileOption() = %v, want %v", got, test.want)
			}
		})
	}
}

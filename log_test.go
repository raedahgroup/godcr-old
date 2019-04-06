package main

import (
	"reflect"
	"testing"
)

func Test_setLogLevel(t *testing.T) {
	tests := []struct {
		name        string
		subsystemID string
		logLevel    string
	}{
		{
			name:        "weblog info level",
			subsystemID: "WEB",
			logLevel:    "info",
		},
		{
			name:        "weblog warn level",
			subsystemID: "WEB",
			logLevel:    "warn",
		},
		{
			name:        "weblog error level",
			subsystemID: "WEB",
			logLevel:    "error",
		},
		{
			name:        "clilog info level",
			subsystemID: "CLI",
			logLevel:    "info",
		},
		{
			name:        "clilog info level",
			subsystemID: "CLI",
			logLevel:    "warn",
		},
		{
			name:        "clilog info level",
			subsystemID: "CLI",
			logLevel:    "error",
		},
		{
			name:        "nuklog info level",
			subsystemID: "NUK",
			logLevel:    "info",
		},
		{
			name:        "nuklog info level",
			subsystemID: "NUK",
			logLevel:    "warn",
		},
		{
			name:        "nuklog info level",
			subsystemID: "NUK",
			logLevel:    "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setLogLevel(tt.subsystemID, tt.logLevel)
		})
	}
}

func Test_setLogLevels(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
	}{
		{
			name:     "weblog info",
			logLevel: "info",
		},
		{
			name:     "weblog warn",
			logLevel: "warn",
		},
		{
			name:     "weblog error",
			logLevel: "error",
		},
		{
			name:     "nuklog info",
			logLevel: "info",
		},
		{
			name:     "nuklog warn",
			logLevel: "warn",
		},
		{
			name:     "nuklog error",
			logLevel: "error",
		},
		{
			name:     "clilog info",
			logLevel: "info",
		},
		{
			name:     "clilog warn",
			logLevel: "warn",
		},
		{
			name:     "clilog error",
			logLevel: "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setLogLevels(tt.logLevel)
		})
	}
}

func Test_validLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		want     bool
	}{
		{
			name:     "valid log level 1",
			logLevel: "info",
			want:     true,
		},
		{
			name:     "valid log level 2",
			logLevel: "warn",
			want:     true,
		},
		{
			name:     "invalid log level",
			logLevel: "notice",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validLogLevel(tt.logLevel); got != tt.want {
				t.Errorf("validLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_supportedSubsystems(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "subsystems",
			want: []string{"CLI", "GODC", "NUK", "TER", "WEB"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := supportedSubsystems(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("supportedSubsystems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseAndSetDebugLevels(t *testing.T) {
	tests := []struct {
		name       string
		debugLevel string
		wantErr    bool
	}{
		{
			name:       "valid debug level",
			debugLevel: "info",
			wantErr:    false,
		},
		{
			name:       "invalid debug level",
			debugLevel: "notice",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseAndSetDebugLevels(tt.debugLevel); (err != nil) != tt.wantErr {
				t.Errorf("parseAndSetDebugLevels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

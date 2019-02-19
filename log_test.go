package main

import (
	"reflect"
	"testing"
)

func Test_logWriter_Write(t *testing.T) {
	tests := []struct {
		name      string
		logWriter logWriter
		p         []byte
		wantN     int
		wantErr   bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logWriter := logWriter{}
			gotN, err := logWriter.Write(tt.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("logWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("logWriter.Write() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func Test_initLogRotator(t *testing.T) {
	tests := []struct {
		name    string
		logFile string
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initLogRotator(tt.logFile)
		})
	}
}

func Test_setLogLevel(t *testing.T) {
	tests := []struct {
		name        string
		subsystemID string
		logLevel    string
	}{
		{
			name:        "weblog",
			subsystemID: "WEB",
			logLevel:    "info",
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseAndSetDebugLevels(tt.debugLevel); (err != nil) != tt.wantErr {
				t.Errorf("parseAndSetDebugLevels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fatalf(t *testing.T) {
	tests := []struct {
		name   string
		format string
		args   []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fatalf(tt.format, tt.args...)
		})
	}
}

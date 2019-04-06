package weblog

import (
	"testing"

	"github.com/decred/slog"
)

func TestDisableLog(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			DisableLog()
		})
	}
}

func TestUseLogger(t *testing.T) {
	tests := []struct {
		name   string
		logger slog.Logger
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			UseLogger(test.logger)
		})
	}
}

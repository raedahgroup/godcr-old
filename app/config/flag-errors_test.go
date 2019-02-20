package config

import (
	"testing"

	flags "github.com/jessevdk/go-flags"
)

func TestIsFlagErrorType(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		errorType flags.ErrorType
		want      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFlagErrorType(tt.err, tt.errorType); got != tt.want {
				t.Errorf("IsFlagErrorType() = %v, want %v", got, tt.want)
			}
		})
	}
}

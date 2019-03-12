package terlog

import "testing"

func TestLogInfo(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			LogInfo(test.message)
		})
	}
}

func TestLogWarn(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			LogWarn(test.message)
		})
	}
}

func TestLogError(t *testing.T) {
	tests := []struct {
		name    string
		message error
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			LogError(test.message)
		})
	}
}

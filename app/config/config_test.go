package config

import "testing"

func TestLoadConfig(t *testing.T) {
	_, _, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load dcrd config: %s", err)
	}
}

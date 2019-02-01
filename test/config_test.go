package test

import (
	"testing"

	"github.com/raedahgroup/godcr/app/config"
)

// TestLoadConfig makes sure that loading config succeeds.
func TestLoadConfig(t *testing.T) {
	_, _, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load godcr config: %s", err)
	}
}

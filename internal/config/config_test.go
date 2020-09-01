package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	k := "LOG_LEVEL"
	old := os.Getenv(k)
	new := "test"
	os.Setenv(k, new)

	cfg, err := Load()
	if err != nil {
		t.Error(err)
	}
	if cfg.LogLevel != new {
		t.Errorf("log level config '%+2v' must be equal to '%s'", cfg, new)
	}

	os.Setenv(k, old)
}

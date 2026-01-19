package config_test

import (
	"testing"

	"github.com/pro200/go-config"
)

func TestEnv(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Error("env load error:", err)
		return
	}

	strVal := cfg.String("STRING")
	intVal := cfg.Int("INT")
	floatVal := cfg.Float64("FLOAT")

	if strVal != "hello" || intVal != 1234 || floatVal != 12.34 {
		t.Error("Wrong result")
	}
}

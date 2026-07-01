package platform

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var config Config

// Config is the top level configuration for an application.
type Config struct {
	// ServiceName is the sercice name for the OTel SDK.
	ServiceName string
	// HumanLogs enables human-readible logs.
	HumanLogs bool
	// TelemetryEnabled turns on all telemetry.
	TelemetryEnabled bool
}

func (c *Config) apply() {
	if c.ServiceName != "" {
		os.Setenv("OTEL_SERVICE_NAME", c.ServiceName)
	}

	if c.HumanLogs {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

// ReadConfig returns a copy of the application configuration.
func ReadConfig() Config {
	return config
}

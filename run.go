// Package platform provides a platform for web services.
package platform

import "github.com/rs/zerolog/log"

// Run is the entrypoint for an application. Optionally provide a service configuration.
func Run(main func() error, c Config) {
	config = c
	config.apply()

	log.Info().Any("config", config).Msg("Starting")
	defer log.Info().Msg("Stopped")

	if err := main(); err != nil {
		log.Error().Err(err).Send()
	}
}

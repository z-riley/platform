// Package platform provides a platform for web services.
package platform

import (
	"context"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
)

// Run is the entrypoint for an application.
func Run(entrypoint func(context.Context) error, c Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config = c
	config.apply()

	log.Info().Any("config", config).Msg("Starting")
	defer log.Info().Msg("Stopped")

	if err := entrypoint(ctx); err != nil {
		log.Error().Err(err).Send()
	}
}

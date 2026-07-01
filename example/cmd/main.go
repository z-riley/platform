package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/z-riley/platform"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	platform.Run(entrypoint, platform.Config{
		ServiceName:      "ExampleService",
		HumanLogs:        true,
		TelemetryEnabled: false,
	})
}

func entrypoint(ctx context.Context) error {
	conn, err := grpc.NewClient(":4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to create new grpc client: %w", err)
	}
	defer conn.Close()

	if platform.ReadConfig().TelemetryEnabled {
		shutdown, err := setupOTelSDK(ctx, conn)
		if err != nil {
			return err
		}
		defer func() {
			err := shutdown(ctx)
			if err != nil {
				log.Error().Err(err).Send()
			}
		}()
	}

	errCh := make(chan error)
	go func() {
		errCh <- http.ListenAndServe(":8080", newHandler())
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return nil
	}
}

func newHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg("/test called")
		fmt.Fprint(w, "hello test")
	})

	handler := otelhttp.NewHandler(mux, "/") // all endpoints
	handler = loggingMiddleware(handler)
	handler = metricMiddleware(handler)
	return handler
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().Str("path", r.URL.Path).Msg("Endpoint called")

		h.ServeHTTP(w, r)
	})
}

func metricMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meter := otel.Meter("testmeter")
		c, err := meter.Int64Counter("testcounter")
		if err != nil {
			panic(err)
		}
		c.Add(r.Context(), 1) // temporary test

		h.ServeHTTP(w, r)
	})
}

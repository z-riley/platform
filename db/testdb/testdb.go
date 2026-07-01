// Package testdb provides database test infrastructure.
package testdb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func New(t *testing.T) string {
	container, err := postgres.Run(t.Context(), "postgres:18.3",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, container)

	url, err := container.ConnectionString(t.Context())
	if err != nil {
		require.NoError(t, err)
	}

	return url
}

// Integration marks a test as an integration test.
func Integration(t *testing.T) {
	const envVar = "INTEGRATION"
	if os.Getenv(envVar) == "" {
		t.Skipf("%s variable not set; skipping", envVar)
	}
}

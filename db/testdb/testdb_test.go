package testdb

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/z-riley/platform/db"
)

func TestNew(t *testing.T) {
	Integration(t)

	url1 := New(t)
	url2 := New(t)

	_, err := db.Connect(t.Context(), url1)
	require.NoError(t, err)

	_, err = db.Connect(t.Context(), url2)
	require.NoError(t, err)
}

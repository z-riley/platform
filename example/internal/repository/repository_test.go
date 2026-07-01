package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/z-riley/platform/db"
	"github.com/z-riley/platform/db/testdb"
)

func TestUpdateSomething(t *testing.T) {
	testdb.Integration(t)

	url := testdb.New(t)

	db, err := db.Connect(t.Context(), url)
	require.NoError(t, err)

	r := Repository{
		db: db,
	}

	err = r.UpdateSomething(t.Context())
	require.NoError(t, err)
}

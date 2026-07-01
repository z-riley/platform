package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/z-riley/platform/db/testdb"
)

func TestRunMigrations(t *testing.T) {
	testdb.Integration(t)

	// Create a migration script.
	dir := t.TempDir()
	f, err := os.Create(filepath.Join(dir, "1_test.sql"))
	require.NoError(t, err)
	_, err = f.WriteString("create table test (id int primary key);")
	require.NoError(t, err)

	// Run the migration.
	url := testdb.New(t)
	err = RunMigrations(t.Context(), url, dir)
	require.NoError(t, err)

	// Check the migration ran.
	conn, err := pgx.Connect(t.Context(), url)
	require.NoError(t, err)
	_, err = conn.Exec(t.Context(), "select * from test")
	require.NoError(t, err)
}

func TestReadMigrationScripts(t *testing.T) {
	for _, tt := range []struct {
		name      string
		fileName  string
		contents  string
		expect    []*migrate.Migration
		expectErr error
	}{
		{
			name:     "success",
			fileName: "1_test.sql",
			contents: "test migration",
			expect: []*migrate.Migration{
				{
					Sequence: 1,
					Name:     "1_test.sql",
					UpSQL:    "test migration",
				},
			},
		},
		{
			name:      "no underscore",
			fileName:  "01test.sql",
			expectErr: errBadMigrationFileName,
		},
		{
			name:      "no number",
			fileName:  "no_number.sql",
			expectErr: errBadMigrationFileName,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			f, err := os.Create(filepath.Join(dir, tt.fileName))
			require.NoError(t, err)
			_, err = f.WriteString(tt.contents)
			require.NoError(t, err)

			actual, err := readMigrationScripts(dir)
			require.ErrorIs(t, err, tt.expectErr)
			if tt.expectErr == nil {
				assert.Equal(t, tt.expect, actual)
			}
		})
	}
}

package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
)

// RunMigrations runs every new migration.
func RunMigrations(ctx context.Context, url, path string) error {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	migrator, err := migrate.NewMigrator(ctx, conn, "schema_version")
	if err != nil {
		return err
	}

	migrator.Migrations, err = readMigrationScripts(path)
	if err != nil {
		return err
	}

	return migrator.Migrate(ctx)
}

var errBadMigrationFileName = errors.New("migration filename must be in the format '1_desc.sql'")

func readMigrationScripts(path string) ([]*migrate.Migration, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var migrations []*migrate.Migration
	for _, file := range files {
		if file.IsDir() || !strings.Contains(file.Name(), ".sql") {
			continue
		}

		before, _, ok := strings.Cut(file.Name(), "_")
		seq, err := strconv.Atoi(before)
		if !ok || err != nil {
			return nil, fmt.Errorf("%w (%s)", errBadMigrationFileName, file.Name())
		}

		b, err := os.ReadFile(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, &migrate.Migration{
			Sequence: int32(seq),
			Name:     file.Name(),
			UpSQL:    string(b),
		})
	}

	return migrations, nil
}

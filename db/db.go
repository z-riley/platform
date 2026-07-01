package db

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// Connect connects to a database.
func Connect(ctx context.Context, url string) (*sqlx.DB, error) {
	return sqlx.ConnectContext(ctx, "pgx", url)
}

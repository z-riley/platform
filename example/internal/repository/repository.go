// Package repository provides a storage layer.
package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db DB
}

func (r *Repository) Begin(ctx context.Context) (*Tx, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	return &Tx{tx: tx}, err
}

func (s *Repository) UpdateSomething(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `UPDATE something SET name = "test" WHERE something_id = 1;`)
	return err
}

func UpdateSomething(ctx context.Context, tx sqlx.Tx) error {
	r := tx.QueryRowContext(ctx, "")
	r.Scan()

	_, err := tx.ExecContext(ctx, `UPDATE something SET name = "test" WHERE something_id = 1;`)
	return err
}

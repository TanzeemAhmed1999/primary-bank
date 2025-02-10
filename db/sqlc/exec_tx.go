package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// execTx executes a func in a db transaction
func (s *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if tErr := tx.Rollback(ctx); tErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, tErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

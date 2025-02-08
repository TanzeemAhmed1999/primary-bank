package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all the functions to execute queries and transactions
type Store struct {
	*Queries
	db *pgxpool.Pool
}

// NewStore creates a new Store
func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

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

// TransferTx performs a money transfer from one account to other.
// It created a transfer record, add account entries, and update accounts balance within a single db
func (s *Store) TransferTx(ctx context.Context, args CreateTransferParams) (TransferTxResult, error) {
	var retval TransferTxResult

	err := s.execTx(ctx, func(queries *Queries) error {
		var err error

		retval.Transfer, err = s.Queries.CreateTransfer(ctx, args)
		if err != nil {
			return err
		}

		retval.FromEntry, err = s.Queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		retval.ToEntry, err = s.Queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		retval.FromAccount, err = s.Queries.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     args.FromAccountID,
			Amount: -args.Amount,
		})
		if err != nil {
			return err
		}

		retval.ToAccount, err = s.Queries.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     args.ToAccountID,
			Amount: args.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return retval, err
	}

	return retval, nil
}

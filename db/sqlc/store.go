package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxRetries = 5
	retryDelay = 100 * time.Millisecond
)

// Store provides all the functions to execute queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, args CreateTransferParams) (TransferTxResult, error)
}

// Store provides all the functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

// NewStore creates a new Store
func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// TransferTx performs a money transfer from one account to other.
// It created a transfer record, add account entries, and update accounts balance within a single db
func (s *SQLStore) TransferTx(ctx context.Context, args CreateTransferParams) (TransferTxResult, error) {
	var retval TransferTxResult
	var err error

	for i := 0; i < maxRetries; i++ {
		err = s.execTx(ctx, func(queries *Queries) error {
			var txErr error

			retval.Transfer, txErr = queries.CreateTransfer(ctx, args)
			if txErr != nil {
				return txErr
			}

			retval.FromEntry, txErr = queries.CreateEntry(ctx, CreateEntryParams{
				AccountID: args.FromAccountID,
				Amount:    -args.Amount,
			})
			if txErr != nil {
				return txErr
			}

			retval.ToEntry, txErr = queries.CreateEntry(ctx, CreateEntryParams{
				AccountID: args.ToAccountID,
				Amount:    args.Amount,
			})
			if txErr != nil {
				return txErr
			}

			if args.FromAccountID < args.ToAccountID {
				retval.FromAccount, retval.ToAccount, txErr = addMoney(ctx, queries, args.FromAccountID, args.ToAccountID, -args.Amount, args.Amount)
			} else {
				retval.FromAccount, retval.ToAccount, txErr = addMoney(ctx, queries, args.ToAccountID, args.FromAccountID, args.Amount, -args.Amount)
			}
			if txErr != nil {
				return txErr
			}

			return nil
		})

		if err == nil {
			return retval, nil
		}

		if !isRetryableError(err) {
			break
		}

		time.Sleep(retryDelay)
	}

	return retval, err
}

func isRetryableError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "40P01"
	}
	return false
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	accountID2 int64,
	amount1 int64,
	amount2 int64,
) (acc1 Account, acc2 Account, err error) {
	acc1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	acc2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return acc1, acc2, err
}

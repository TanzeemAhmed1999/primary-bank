package db

import (
	"context"

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

		if args.FromAccountID < args.ToAccountID {
			retval.FromAccount, retval.ToAccount, err = addMoney(ctx, queries, args.FromAccountID, args.ToAccountID, -args.Amount, args.Amount)
		} else {
			retval.FromAccount, retval.ToAccount, err = addMoney(ctx, queries, args.ToAccountID, args.FromAccountID, args.Amount, -args.Amount)
		}
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

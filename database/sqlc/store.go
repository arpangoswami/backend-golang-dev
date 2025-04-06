package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides a interface to implement transactions
type Store struct {
	*Queries
	database *sql.DB
}

// NewStore returns a instance of Store object
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries:  New(db),
		database: db,
	}
}

func (store *Store) executeTransaction(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error %w; rollback error failed: %w", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTransactionParams struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type TransferTransactionResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTransaction performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts' balance within a single db txn
func (store *Store) TransferTransaction(ctx context.Context, arg TransferTransactionParams) (TransferTransactionResult, error) {
	var result TransferTransactionResult
	err := store.executeTransaction(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		// TODO: update accounts balance
		return nil
	})
	return result, err
}

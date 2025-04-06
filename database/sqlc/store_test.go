package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStore_TransferTransaction(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := float64(5)

	errs := make(chan error)
	results := make(chan TransferTransactionResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTransaction(context.Background(), TransferTransactionParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check whether errors are empty
	for i := 0; i < n; i++ {
		err := <-errs
		assert.Nil(t, err)
		transactionResult := <-results
		assert.NotNil(t, transactionResult)
		assert.NotEmpty(t, transactionResult)

		transferResult := transactionResult.Transfer

		assert.Equal(t, account1.ID, transferResult.FromAccountID)
		assert.Equal(t, account2.ID, transferResult.ToAccountID)
		assert.Equal(t, amount, transferResult.Amount)
		assert.NotZero(t, transferResult.ID)
		assert.NotZero(t, transferResult.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transferResult.ID)
		assert.NoError(t, err)

		fromEntry := transactionResult.FromEntry
		assert.NotEmpty(t, fromEntry)
		assert.Equal(t, account1.ID, fromEntry.AccountID)
		assert.Equal(t, -amount, fromEntry.Amount)
		assert.NotZero(t, fromEntry.ID)
		assert.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		assert.NoError(t, err)

		toEntry := transactionResult.ToEntry
		assert.NotEmpty(t, toEntry)
		assert.Equal(t, account2.ID, toEntry.AccountID)
		assert.Equal(t, amount, toEntry.Amount)
		assert.NotZero(t, toEntry.ID)
		assert.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		assert.NoError(t, err)

		// TODO: check accounts balance. Tackle it with lock of store transaction
	}
}

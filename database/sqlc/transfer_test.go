package db

import (
	"context"
	"database/sql"
	"github.com/arpangoswami/backend-golang-dev/util"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T) Transfer {
	t.Helper()
	listAccountArg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), listAccountArg)
	assert.NotEmpty(t, accounts)
	assert.Nil(t, err)
	n := len(accounts)
	randIdx := rand.Intn(n)
	fromAccountId := accounts[randIdx].ID
	fromBalance := accounts[randIdx].Balance
	money := util.RandomMoney(int64(accounts[randIdx].Balance))
	randIdx = rand.Intn(n)
	toAccountId := accounts[randIdx].ID
	toBalance := accounts[randIdx].Balance
	createTransferArgs := CreateTransferParams{
		FromAccountID: fromAccountId,
		ToAccountID:   toAccountId,
		Amount:        money,
	}
	transfer, err :=
		testQueries.CreateTransfer(context.Background(), createTransferArgs)
	assert.Nil(t, err)
	assert.NotEmpty(t, transfer)
	assert.Equal(t, fromAccountId, transfer.FromAccountID)
	assert.Equal(t, toAccountId, transfer.ToAccountID)
	assert.Equal(t, money, transfer.Amount)
	_, err = testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      fromAccountId,
		Balance: fromBalance - money,
	})
	assert.Nil(t, err)
	_, err = testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      toAccountId,
		Balance: toBalance + money,
	})
	return transfer
}

func TestQueries_CreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestQueries_DeleteTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	cleanupTransfer(t, transfer1.ID)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, transfer2)
}

func TestQueries_GetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	assert.NoError(t, err)
	assert.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	assert.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	assert.Equal(t, transfer1.Amount, transfer2.Amount)
	assert.WithinDurationf(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second,
		"createdAt1 and createdAt2 should have max difference of 1s")
	cleanupTransfer(t, transfer1.ID)
}

type pairFromToAccountId struct {
	FromAccountId int64
	ToAccountId   int64
}

func TestQueries_ListTransfers(t *testing.T) {
	var cleanupList []int64
	var transferPairs []pairFromToAccountId
	for i := 0; i < 10; i++ {
		transfer := createRandomTransfer(t)
		transferPairs = append(transferPairs, pairFromToAccountId{
			FromAccountId: transfer.FromAccountID,
			ToAccountId:   transfer.ToAccountID,
		})
		cleanupList = append(cleanupList, transfer.ID)
	}
	for _, transfer := range transferPairs {
		listTransferArgs := ListTransfersParams{
			FromAccountID: transfer.FromAccountId,
			ToAccountID:   transfer.ToAccountId,
			Limit:         5,
			Offset:        5,
		}
		transfers, err := testQueries.ListTransfers(context.Background(), listTransferArgs)
		assert.NoError(t, err)
		if len(transfers) > 0 {
			for _, listTransferResult := range transfers {
				assert.NotEmpty(t, listTransferResult)
			}
		}
		assert.LessOrEqual(t, len(transfers), 5)
	}

	cleanupTransfers(t, cleanupList)
}

func undoTransfer(t *testing.T, transferId int64) {
	t.Helper()
	transfer, err := testQueries.GetTransfer(context.Background(), transferId)
	assert.NoError(t, err)
	assert.Equal(t, transferId, transfer.ID)
	amount := transfer.Amount
	fromAccount, err := testQueries.GetAccount(context.Background(), transfer.FromAccountID)
	assert.NoError(t, err)
	toAccount, err := testQueries.GetAccount(context.Background(), transfer.ToAccountID)
	_, err = testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      fromAccount.ID,
		Balance: fromAccount.Balance + amount,
	})
	assert.NoError(t, err)
	_, err = testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      toAccount.ID,
		Balance: toAccount.Balance - amount,
	})
	assert.NoError(t, err)
}

func cleanupTransfer(t *testing.T, transferId int64) {
	t.Helper()
	undoTransfer(t, transferId)
	err := testQueries.DeleteTransfer(context.Background(), transferId)
	assert.NoError(t, err)
}

func cleanupTransfers(t *testing.T, transferIds []int64) {
	t.Helper()
	var err error
	for _, id := range transferIds {
		cleanupTransfer(t, id)
		assert.NoError(t, err)
	}
}

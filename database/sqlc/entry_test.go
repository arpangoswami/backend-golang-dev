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

// Entry represents a deposit or a withdrawal. For a withdrawal check if withdrawal amount < balance
func createRandomEntry(t *testing.T) Entry {
	t.Helper()
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	assert.NotEmpty(t, accounts)
	assert.Nil(t, err)
	n := len(accounts)
	var money float64
	fromRandIdx := rand.Intn(n)
	fromAccountId := accounts[fromRandIdx].ID
	toRandIdx := rand.Intn(n)
	toAccountId := accounts[toRandIdx].ID
	money = util.RandomMoney(int64(accounts[fromRandIdx].Balance))
	_, err = testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      fromAccountId,
		Balance: accounts[fromRandIdx].Balance - money,
	})
	assert.NoError(t, err)
	_, err = testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      toAccountId,
		Balance: accounts[toRandIdx].Balance + money,
	})
	entry, err := testQueries.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: fromAccountId,
		Amount:    -money,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fromAccountId, entry.AccountID)
	assert.Equal(t, -money, entry.Amount)
	assert.Nil(t, err)
	receivedEntry, err := testQueries.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: toAccountId,
		Amount:    money,
	})
	assert.NoError(t, err)
	assert.Equal(t, toAccountId, receivedEntry.AccountID)
	assert.Equal(t, money, receivedEntry.Amount)
	return entry
}

func TestQueries_CreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestQueries_GetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	assert.NoError(t, err)
	assert.Equal(t, entry1.AccountID, entry2.AccountID)
	assert.Equal(t, entry1.Amount, entry2.Amount)
	assert.WithinDurationf(t, entry1.CreatedAt, entry2.CreatedAt, time.Second,
		"createdAt1 and createdAt2 should have max difference of 1s")
	cleanUpEntry(t, entry1.ID)
}

func TestQueries_DeleteEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	cleanUpEntry(t, entry1.ID)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, entry2)
}

func TestQueries_ListEntries(t *testing.T) {
	var cleanupList []int64
	accountIDs := make(map[int64]bool)
	for i := 0; i < 10; i++ {
		entry := createRandomEntry(t)
		accountIDs[entry.AccountID] = true
		cleanupList = append(cleanupList, entry.ID)
	}
	arg := ListEntriesParams{
		Column1: getUniqueAccountIDs(t, accountIDs),
		Limit:   5,
		Offset:  5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	assert.NoError(t, err)
	assert.Len(t, entries, 5)

	for _, entry := range entries {
		assert.NotEmpty(t, entry)
	}

	cleanUpEntries(t, cleanupList)
}

func undoEntry(t *testing.T, entryId int64) {
	entry, err := testQueries.GetEntry(context.Background(), entryId)
	assert.NoError(t, err)
	amount := entry.Amount
	previousAmount, err := testQueries.GetAccount(context.Background(), entry.AccountID)
	assert.NoError(t, err)
	testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      entryId,
		Balance: previousAmount.Balance - amount,
	})
}

func cleanUpEntry(t *testing.T, id int64) {
	t.Helper()
	undoEntry(t, id)
	err := testQueries.DeleteEntry(context.Background(), id)
	assert.NoError(t, err)
}

func cleanUpEntries(t *testing.T, ids []int64) {
	t.Helper()
	var err error
	for _, id := range ids {
		cleanUpEntry(t, id)
		assert.NoError(t, err)
	}
}

func getUniqueAccountIDs(t *testing.T, accountIDs map[int64]bool) []int64 {
	t.Helper()
	keys := make([]int64, len(accountIDs))
	ptr := 0
	for k := range accountIDs {
		keys[ptr] = k
		ptr += 1
	}
	return keys
}

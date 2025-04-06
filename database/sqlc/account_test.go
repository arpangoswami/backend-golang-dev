package db

import (
	"context"
	"database/sql"
	"github.com/arpangoswami/backend-golang-dev/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account {
	t.Helper()
	currencyCountryCode := util.RandomCurrencyCodeCountryCode()
	arg := CreateAccountParams{
		Owner:       util.RandomOwner(),
		Balance:     util.RandomMoney(100000),
		Currency:    currencyCountryCode.CurrencyCode,
		CountryCode: currencyCountryCode.CountryCode,
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	assert.NotNil(t, account)
	assert.NoError(t, err)
	assert.Equal(t, arg.Owner, account.Owner)
	assert.Equal(t, arg.Balance, account.Balance)
	assert.Equal(t, arg.Currency, account.Currency)
	assert.Equal(t, arg.CountryCode, account.CountryCode)
	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)
	return account
}

// Not using the cleanup inside CreateAccount to validate that records are being inserted in the database
func TestQueries_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestQueries_GetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	assert.Equal(t, account1.CountryCode, account2.CountryCode)
	assert.WithinDurationf(t, account1.CreatedAt, account2.CreatedAt, time.Second,
		"createdAt1 and createdAt2 should have max difference of 1s")
	cleanUpAccount(t, account1.ID)
}

func TestQueries_UpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(100000),
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.Equal(t, account1.ID, account2.ID)
	assert.NotEqual(t, account1.Balance, account2.Balance)
	assert.Equal(t, arg.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	assert.Equal(t, account1.CountryCode, account2.CountryCode)
	assert.WithinDurationf(t, account1.CreatedAt, account2.CreatedAt, time.Second,
		"createdAt1 and createdAt2 should have max difference of 1s")
	cleanUpAccount(t, account1.ID)
}

func TestQueries_DeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	cleanUpAccount(t, account1.ID)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, account2)
}

func TestQueries_ListAccounts(t *testing.T) {
	var cleanupList []int64
	for i := 0; i < 10; i++ {
		account := createRandomAccount(t)
		cleanupList = append(cleanupList, account.ID)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	assert.NoError(t, err)
	assert.Len(t, accounts, 5)

	for _, account := range accounts {
		assert.NotEmpty(t, account)
	}

	cleanUpAccounts(t, cleanupList)
}

func cleanUpAccount(t *testing.T, id int64) {
	t.Helper()
	err := testQueries.DeleteAccount(context.Background(), id)
	assert.NoError(t, err)
}

func cleanUpAccounts(t *testing.T, ids []int64) {
	t.Helper()
	var err error
	for _, id := range ids {
		err = testQueries.DeleteAccount(context.Background(), id)
		assert.NoError(t, err)
	}
}

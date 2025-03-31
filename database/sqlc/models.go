// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"database/sql"
	"time"
)

type Account struct {
	ID          int64         `json:"id"`
	Owner       string        `json:"owner"`
	Balance     float64       `json:"balance"`
	Currency    string        `json:"currency"`
	CreatedAt   time.Time     `json:"created_at"`
	CountryCode sql.NullInt32 `json:"country_code"`
}

type Entry struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	// Can be both negative and positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Transfer struct {
	ID            int64 `json:"id"`
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	// Must be positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

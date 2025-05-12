package database

import "time"

type AccountModel struct {
	ID        int
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionModel struct {
	ID              int
	AccountID       int
	TransactionDate time.Time
	Amount          float64
	TransactionType string
	CreatedAt       time.Time
}

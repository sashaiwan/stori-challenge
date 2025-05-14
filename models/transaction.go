package models

import (
	"stori-challenge/database"
	"time"
)

type TransactionType int

const (
	Debit TransactionType = iota
	Credit
)

type Transaction struct {
	Id     int
	Date   string
	Amount float64
	Type   TransactionType
}

func (t Transaction) ToDatabaseModel(accountID int) database.TransactionModel {
	// We already ensure date integrity
	date, _ := time.Parse("2006-01-02", t.Date)

	transactionType := "credit"
	if t.Type == Debit {
		transactionType = "debit"
	}

	return database.TransactionModel{
		ID:              t.Id,
		AccountID:       accountID,
		TransactionDate: date,
		Amount:          t.Amount,
		TransactionType: transactionType,
	}
}

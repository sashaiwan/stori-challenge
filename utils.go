package main

import (
	"stori-challenge/database"
	"stori-challenge/models"
)

// This can be generic
func getTransactionsModels(accountID int, transactions []models.Transaction) []database.TransactionModel {
	models := make([]database.TransactionModel, len(transactions))

	for i, t := range transactions {
		models[i] = t.ToDatabaseModel(accountID)
	}

	return models
}

package utils

import (
	"stori-challenge/database"
	"stori-challenge/models"
)

// This can be generic - DON'T DO IT NOW
func GetTransactionsModels(accountID int, transactions []models.Transaction) []database.TransactionModel {
	models := make([]database.TransactionModel, len(transactions))

	for i, t := range transactions {
		models[i] = t.ToDatabaseModel(accountID)
	}

	return models
}

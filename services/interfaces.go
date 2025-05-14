package services

import "stori-challenge/models"

type EmailSender interface {
	SendEmail(stat TransactionStats, recipient string) error
}

type TransactionProcessor interface {
	GetTransactionStats(transactions []models.Transaction) TransactionStats
}

type FilesProcessor interface {
	ProcessCSV(filepath string) ([]models.Transaction, error)
}

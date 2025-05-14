package server

import (
	"stori-challenge/database"
	"stori-challenge/services"
)

type Server struct {
	emailSender          services.EmailSender
	transactionProcessor services.TransactionProcessor
	csvProcessor         services.FilesProcessor
	db                   *database.DB
}

func NewServer(emailSender services.EmailSender, transactionProcessor services.TransactionProcessor, csvProcessor services.FilesProcessor, db *database.DB) *Server {
	return &Server{
		emailSender:          emailSender,
		transactionProcessor: transactionProcessor,
		csvProcessor:         csvProcessor,
		db:                   db,
	}
}

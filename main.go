package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"stori-challenge/database"
	"stori-challenge/server"
	"stori-challenge/services"

	"github.com/joho/godotenv"
)

func main() {
	// TODO:
	// Naive approach:
	// 1. Create a CSV reader that returns struct with the transfer data (or something similar)
	// 2. Create the function to persist on a DB
	// 3. Create the email sender function/service

	// Event driven approach:
	// 1. Define the representation of the events
	// 2. Create the event store (or emulate it with a classic SQL, research)
	// 3. Create the commands to persist the events
	// 4. Create the query for retrieve the data
	// 5. Create the CSV parser
	// 6. Create the email service

	if !isDocker() {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, using system environment")
		}
	} else {
		log.Println("Running in Docker, skipping .env file")
	}

	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	} else {
		fmt.Println("DB initialized and ready")
		defer db.Close()
	}

	// TODO: create a config mechanism to unify env variables handling
	emailService := &services.EmailService{}
	transactionsService := &services.TransactionService{}
	filesService := &services.FilesService{}
	server := server.NewServer(emailService, transactionsService, filesService, db)

	http.HandleFunc("/transactions/stats", server.TransactionStatsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting HTTP server on post %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("HTTP server error: ", err)
	}
}

func isDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

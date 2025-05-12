package main

import (
	"fmt"
	"os"
	"stori-challenge/database"

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

	godotenv.Load()

	// TODO: make the file path setting more reliable
	csvPath := os.Getenv("CSV_FILE_PATH")
	if csvPath == "" {
		// Should use filePath.Join
		csvPath = "./data/txns.csv"
	}

	recipientMail := os.Getenv("RECIPIENT_EMAIL")
	if recipientMail == "" {
		fmt.Fprintf(os.Stderr, "A recipient mail must be provided")
		os.Exit(1)
	}

	db, err := database.NewDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
	} else {
		fmt.Println("DB initialized and ready")
		defer db.Close()
	}

	account, err := db.GetOrCreateAccount(recipientMail)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching or creating account: %v", err)
	}

	transactions, err := processCSV(csvPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully processed %d transactions\n", len(transactions))

	result := db.SaveTransactions(account.ID, getTransactionsModels(account.ID, transactions))
	if result != nil {
		fmt.Fprintf(os.Stderr, "Error saving transactions: %v\n", result)
	} else {
		fmt.Printf("Successfully saved %d transactions\n", len(transactions))
	}

	transactionStats := getTransactionStats(transactions)
	fmt.Println(transactionStats)

	mailErr := sendEmail(transactionStats, recipientMail)
	if mailErr != nil {
		fmt.Fprintf(os.Stderr, "Error sending email: %v\n", mailErr)
		os.Exit(1)
	}
}

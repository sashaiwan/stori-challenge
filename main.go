package main

import (
	"fmt"
	"os"
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

	filepath := "./txns.csv"
	transactions, err := processCSV(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %d transactions\n", len(transactions))
}

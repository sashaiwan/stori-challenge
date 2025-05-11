package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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

func processCSV(filepath string) ([]Transaction, error) {

	lines, err := readAllLines(filepath)
	if err != nil {
		return nil, err
	}

	// TODO: find a better way to validate this
	header := strings.Split(lines[0], ",")
	expectedHeader := []string{"id", "date", "transaction"}
	for i := range header {
		header[i] = strings.ToLower(header[i])
	}
	if !validateHeader(header, expectedHeader) {
		return nil, fmt.Errorf("Invalid CSV header. Expected: %v, Actual: %v",
			expectedHeader, header)
	}

	var transactions []Transaction

	// We already process and validate the header
	for _, line := range lines[1:] {

		fields := strings.Split(line, ",")
		transaction, err := parseTransaction(fields)
		if err != nil {
			fmt.Printf("Invalid transaction record: %v", err)
			continue
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func readAllLines(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Error reading all lines: %v", err)
	}

	dataContent := string(data)
	lines := strings.Split(dataContent, "\n")

	var cleanedLines []string
	for _, line := range lines {
		cleanedLine := strings.Trim(line, "\n\r\t ")
		if cleanedLine != "" {
			cleanedLines = append(cleanedLines, cleanedLine)
		}
	}

	return cleanedLines, nil
}

func validateHeader(actual []string, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}

	for i, h := range actual {
		if h != expected[i] {
			return false
		}
	}

	return true
}

func parseTransaction(fields []string) (Transaction, error) {
	// The ID could be an UUID
	id, err := strconv.Atoi(fields[0])
	if err != nil {
		return Transaction{}, fmt.Errorf("Invalid ID format: %v", err)
	}

	date, err := parseDate(fields[1])

	amountStr := fields[2]
	var transactionType TransactionType
	if strings.HasPrefix(amountStr, "+") {
		transactionType = Credit
	} else if strings.HasPrefix(amountStr, "-") {
		transactionType = Debit
	} else {
		return Transaction{}, fmt.Errorf(
			"Transaction prefix is invalid. Expect '-' or '+' but received: %s", amountStr)
	}

	amount, err := strconv.ParseFloat(amountStr[1:], 64)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid amount: %s", err)
	}

	return Transaction{
		Id:     id,
		Date:   date,
		Amount: amount,
		Type:   transactionType,
	}, nil
}

func parseDate(dateStr string) (string, error) {
	parts := strings.Split(dateStr, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("Invalid date format: %s", dateStr)
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return "", fmt.Errorf("Invalid month: %s", parts[0])
	}

	day, err := strconv.Atoi(parts[1])
	if err != nil || day < 1 || day > 31 {
		return "", fmt.Errorf("Invalid day: %s", parts[1])
	}

	year := time.Now().Year()

	return fmt.Sprintf("%d-%02d-%02d", year, month, day), nil
}

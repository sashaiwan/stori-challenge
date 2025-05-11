package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// processCSV gets a filepath and return an array of validated Transaction struct
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

// readAllLines handles the CSV file and return an array of cleaned lines
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

// validateHeader is a helper function to compare the deep equality of two arrays of string
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

// parseTransaction gets an array of string with the parsed fields
// and returns a Transaction struct.
func parseTransaction(fields []string) (Transaction, error) {
	// This could be an UUID
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

// parseDate gets a raw date string with an arbitrary MM/DD format
// and converts it to ISO 8601 format (YYYY-MM-DD)
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

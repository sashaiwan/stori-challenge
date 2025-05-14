package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"stori-challenge/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (s *Server) TransactionStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}
	fmt.Println(r.Header.Get("Content-Type"))
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RequestSchema
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println(req)

	csvData, err := downloadFromS3(req.BucketName, req.ObjectKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to download file: %v", err), http.StatusInternalServerError)
		return
	}

	tempFile, err := os.CreateTemp("", "transactions-*.csv")
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(csvData); err != nil {
		http.Error(w, "Failed to write temp file", http.StatusInternalServerError)
		return
	}
	tempFile.Close()

	account, err := s.db.GetOrCreateAccount(req.Email)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching or creating account: %v", err)
	}
	// TODO: change csvProcessor to receive a reader instead?
	// Creating the temp file allows some retrying mechanism
	// What could happen with a huge file?
	transactions, err := s.csvProcessor.ProcessCSV(tempFile.Name())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing file: %v\n", err)
	}
	log.Printf("Successfully processed %d transactions\n", len(transactions))

	result := s.db.SaveTransactions(account.ID, utils.GetTransactionsModels(account.ID, transactions))
	if result != nil {
		// When failing we can retry since the temp file stills there
		fmt.Fprintf(os.Stderr, "Error saving transactions: %v\n", result)
	} else {
		fmt.Printf("Successfully saved %d transactions\n", len(transactions))
	}

	transactionStats := s.transactionProcessor.GetTransactionStats(transactions)
	fmt.Println(transactionStats)

	mailErr := s.emailSender.SendEmail(transactionStats, req.Email)
	if mailErr != nil {
		fmt.Fprintf(os.Stderr, "Error sending email: %v\n", mailErr)
		http.Error(w, "Failed to send mail summary", http.StatusInternalServerError)
	}

	response := ResponseSchema{
		Status:  "success",
		Message: "Transactions processed and email sent successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "Failed to create response", http.StatusInternalServerError)
		return
	}
}

// TODO: move to a service
func downloadFromS3(bucket, key string) ([]byte, error) {
	log.Printf("Attempting to download from S3 - Bucket: '%s', Key: '%s'", bucket, key)

	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

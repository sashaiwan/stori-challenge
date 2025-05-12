package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func NewDB() (*DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Default values for development
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPassword == "" {
		dbPassword = "password"
	}
	if dbName == "" {
		dbName = "transactions"
	}

	connectionStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	conn, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	db := &DB{conn: conn}

	if err := db.InitSchema(); err != nil {
		return nil, fmt.Errorf("error initializing schema: %v", err)
	}

	return db, nil
}

func (db *DB) InitSchema() error {
	schema := `
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER PRIMARY KEY,
			account_id INTEGER REFERENCES accounts(id),
			transaction_date DATE NOT NULL,
			amount DECIMAL(10,2) NOT NULL,
			transaction_type VARCHAR(10) NOT NULL CHECK (transaction_type IN ('credit', 'debit')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.conn.Exec(schema)
	return err
}

// TODO: explore common Go patterns for db modeling
func (db *DB) GetOrCreateAccount(email string) (*AccountModel, error) {
	var account AccountModel

	err := db.conn.QueryRow("SELECT id, email, created_at, updated_at FROM accounts WHERE email = $1", email).
		Scan(&account.ID, &account.Email, &account.CreatedAt, &account.UpdatedAt)

	if err == sql.ErrNoRows {
		err = db.conn.QueryRow(
			"INSERT INTO accounts (email) VALUES ($1) RETURNING id, email, created_at, updated_at",
			email,
		).Scan(&account.ID, &account.Email, &account.CreatedAt, &account.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("error creating account: %v", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("error fetching account: %v", err)
	}

	return &account, nil
}

func (db *DB) SaveTransactions(accountID int, transactions []TransactionModel) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        INSERT INTO transactions (id, account_id, transaction_date, amount, transaction_type)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE SET
            amount = EXCLUDED.amount,
            transaction_type = EXCLUDED.transaction_type
    `)
	if err != nil {
		return fmt.Errorf("error preparing transaction: %v", err)
	}
	defer stmt.Close()

	for _, trans := range transactions {
		_, err = stmt.Exec(trans.ID, accountID, trans.TransactionDate, trans.Amount, trans.TransactionType)
		if err != nil {
			return fmt.Errorf("error inserting transaction %d: %v", trans.ID, err)
		}
	}

	return tx.Commit()
}

func (db *DB) Close() error {
	return db.conn.Close()
}

package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"payment-service/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	pq "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var (
	testDB *sql.DB
	repo   domain.TransactionRepository
)

const (
	dbHost = "localhost"
	dbPort = "5432"
	dbUser = "user_payment"
	dbPass = "pass_payment"
	dbName = "db_payment_test" // Use a separate database for testing
)

func TestMain(m *testing.M) {
	err := setupTestDB()
	if err != nil {
		log.Fatalf("Failed to set up test database: %v", err)
	}
	defer teardownTestDB()

	repo = NewPostgresRepo(testDB)
	os.Exit(m.Run())
}

func setupTestDB() error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass)
	var err error
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres to create test db: %w", err)
	}
	defer testDB.Close()

	// Create test database if it doesn't exist
	_, err = testDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		// If error is "database ... already exists", ignore it
		if !isDuplicateDatabaseError(err) {
			return fmt.Errorf("failed to create test database: %w", err)
		}
	}

	testDB, err = sql.Open("postgres", fmt.Sprintf("%s dbname=%s", connStr, dbName))
	if err != nil {
		return fmt.Errorf("failed to connect to test db: %w", err)
	}

	// Apply migrations or create tables
	// For simplicity, we'll create tables directly here.
	createWalletsTableSQL := `
	CREATE TABLE IF NOT EXISTS wallets (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL UNIQUE,
		balance BIGINT NOT NULL DEFAULT 0,
		version INT NOT NULL DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	createTransactionsTableSQL := `
	CREATE TABLE IF NOT EXISTS transactions (
		id VARCHAR(36) PRIMARY KEY,
		reference_id VARCHAR(255) NOT NULL UNIQUE,
		sender_id VARCHAR(36) NOT NULL,
		receiver_id VARCHAR(36) NOT NULL,
		amount BIGINT NOT NULL,
		status VARCHAR(50) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = testDB.Exec(createWalletsTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create wallets table: %w", err)
	}
	_, err = testDB.Exec(createTransactionsTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	return nil
}

func teardownTestDB() {
	if testDB != nil {
		testDB.Close()
	}

	// Drop the test database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to connect to postgres to drop test db: %v", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	if err != nil {
		log.Printf("Failed to drop test database: %v", err)
	}
}

func isDuplicateDatabaseError(err error) bool {
	pqErr, ok := err.(*pq.Error)
	return ok && pqErr.Code.Name() == "duplicate_database"
}

func clearTables() error {
	_, err := testDB.Exec("DELETE FROM transactions")
	if err != nil {
		return err
	}
	_, err = testDB.Exec("DELETE FROM wallets")
	return err
}

func TestPostgresRepo_TopUpWallet(t *testing.T) {
	require.NoError(t, clearTables())

	userID := uuid.New().String()
	walletID := uuid.New().String()
	initialBalance := int64(1000)
	topUpAmount := int64(500)

	// Insert initial wallet
	_, err := testDB.Exec(`INSERT INTO wallets (id, user_id, balance, version, created_at, updated_at) 
                           VALUES ($1, $2, $3, $4, $5, $6)`,
		walletID, userID, initialBalance, 0, time.Now(), time.Now())
	require.NoError(t, err)

	tx, err := repo.BeginTx()
	require.NoError(t, err)
	defer repo.RollbackTx(tx)

	err = repo.TopUpWallet(tx, userID, topUpAmount)
	require.NoError(t, err)

	err = repo.CommitTx(tx)
	require.NoError(t, err)

	// Verify balance
	var finalBalance int64
	err = testDB.QueryRow(`SELECT balance FROM wallets WHERE user_id = $1`, userID).Scan(&finalBalance)
	require.NoError(t, err)
	require.Equal(t, initialBalance+topUpAmount, finalBalance)

	// Test case: User not found - TopUpWallet should return an error from GetWalletForUpdate
	require.NoError(t, clearTables()) // Clear for next test case
	tx2, err := repo.BeginTx()
	require.NoError(t, err)
	defer repo.RollbackTx(tx2)

	nonExistentUserID := uuid.New().String()
	err = repo.TopUpWallet(tx2, nonExistentUserID, topUpAmount)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no rows in result set") // Expecting an error from GetWalletForUpdate
	repo.RollbackTx(tx2)                                      // Rollback explicitly since no commit will happen
}

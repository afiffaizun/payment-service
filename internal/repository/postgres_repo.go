package repository

import (
	"database/sql"
	"payment-service/internal/domain"
	"fmt"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) domain.TransactionRepository {
	return &PostgresRepo{db: db}
}


func (r *PostgresRepo) BeginTx() (interface{}) {
	return r.db.Begin()
}

func (r *PostgresRepo) CommitTx(tx interface{}) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}
	return sqlTx.Commit()
}

func (r *PostgresRepo) RollbackTx(tx interface{}) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}
	return sqlTx.Rollback()
}

func (r *PostgresRepo) GetWalletForUpdate(tx interface{}, userID string) (*domain.Wallet, error) {
    sqlTx := tx.(*sql.Tx)
    query := `SELECT id, user_id, balance, version FROM wallets WHERE user_id = $1 FOR UPDATE`
    
    row := sqlTx.QueryRow(query, userID)
    var w domain.Wallet
    err := row.Scan(&w.ID, &w.UserID, &w.Balance, &w.Version)
    if err != nil {
        return nil, err
    }
    return &w, nil
}

func (r *PostgresRepo) UpdateWalletBalance(tx interface{}, walletID string, amount int64) error {
    sqlTx := tx.(*sql.Tx)
    query := `UPDATE wallets SET balance = balance + $1, version = version + 1 WHERE id = $2`
    _, err := sqlTx.Exec(query, amount, walletID)
    return err
}

func (r *PostgresRepo) CreateTransaction(tx interface{}, t *domain.Transaction) error {
    sqlTx := tx.(*sql.Tx)
    query := `INSERT INTO transactions (id, reference_id, sender_id, receiver_id, amount, status, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
    _, err := sqlTx.Exec(query, t.ID, t.ReferenceID, t.SenderID, t.ReceiverID, t.Amount, t.Status, t.CreatedAt)
    return err
}

func (r *PostgresRepo) GetTransactionByRef(refID string) (*domain.Transaction, error) {
    query := `SELECT id, reference_id, sender_id, receiver_id, amount, status, created_at FROM transactions WHERE reference_id = $1`
    row := r.db.QueryRow(query, refID)
    var t domain.Transaction
    err := row.Scan(&t.ID, &t.ReferenceID, &t.SenderID, &t.ReceiverID, &t.Amount, &t.Status, &t.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &t, nil
}


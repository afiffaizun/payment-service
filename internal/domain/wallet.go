package domain

import "time"

type Wallet struct {
	ID        string
	UserID    string
	Balance   int64
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transaction struct {
	ID         string
	Reference  string
	SenderID   string
	ReceiverID string
	Amount     int64
	Status     string
	CreatedAt  time.Time
}

type TransactionRepository interface {
	GetWalletForUpdate(tx interface{}, userID string) (*Wallet, error)
	UpdateWalletBalance(tx interface{}, walletID string, amount int64) error
	CreateTransaction(tx interface{}, transaction *Transaction) error
	GetTransactionByRef(refID string) (*Transaction, error)
	BeginTx() (interface{}, error)
	CommitTx(tx interface{}) error
	RollbackTx(tx interface{}) error
	TopUpWallet(tx interface{}, userID string, amount int64) error
	GetWalletByUserID(userID string) (*Wallet, error)
}

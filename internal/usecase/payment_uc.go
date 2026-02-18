package usecase

import (
	"errors"
	"payment-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount       = errors.New("amount must be greater than zero")
	ErrSameUser            = errors.New("cannot transfer to the same user")
	ErrReferenceExists     = errors.New("reference ID already exists")
)

type PaymentUsecase struct {
	repo domain.TransactionRepository
}

func NewPaymentUsecase(repo domain.TransactionRepository) *PaymentUsecase {
	return &PaymentUsecase{repo: repo}
}

type TransferRequest struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Amount     int64  `json:"amount"`
	Reference  string `json:"reference"`
}

type TransferResponse struct {
	TransactionID string
	Reference     string
	Amount        int64
	Status        string
	CreatedAt     time.Time
}

type TopUpRequest struct {
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
}

type TopUpResponse struct {
	UserID  string
	Amount  int64
	Balance int64
}

type GetWalletResponse struct {
	WalletID  string    `json:"wallet_id"`
	UserID    string    `json:"user_id"`
	Balance   int64     `json:"balance"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *PaymentUsecase) TransferFunds(req TransferRequest) (*TransferResponse, error) {
	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if req.SenderID == req.ReceiverID {
		return nil, ErrSameUser
	}

	existingTx, err := u.repo.GetTransactionByRef(req.Reference)
	if err == nil && existingTx != nil {
		return nil, ErrReferenceExists
	}

	tx, err := u.repo.BeginTx()
	if err != nil {
		return nil, err
	}
	defer u.repo.RollbackTx(tx)

	senderWallet, err := u.repo.GetWalletForUpdate(tx, req.SenderID)
	if err != nil {
		return nil, err
	}

	if senderWallet.Balance < req.Amount {
		return nil, ErrInsufficientBalance
	}

	receiverWallet, err := u.repo.GetWalletForUpdate(tx, req.ReceiverID)
	if err != nil {
		return nil, err
	}

	err = u.repo.UpdateWalletBalance(tx, senderWallet.ID, -req.Amount)
	if err != nil {
		return nil, err
	}

	err = u.repo.UpdateWalletBalance(tx, receiverWallet.ID, req.Amount)
	if err != nil {
		return nil, err
	}

	transaction := &domain.Transaction{
		ID:         uuid.New().String(),
		Reference:  req.Reference,
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Amount:     req.Amount,
		Status:     "completed",
		CreatedAt:  time.Now(),
	}

	err = u.repo.CreateTransaction(tx, transaction)
	if err != nil {
		return nil, err
	}

	err = u.repo.CommitTx(tx)
	if err != nil {
		return nil, err
	}

	return &TransferResponse{
		TransactionID: transaction.ID,
		Reference:     transaction.Reference,
		Amount:        transaction.Amount,
		Status:        transaction.Status,
		CreatedAt:     transaction.CreatedAt,
	}, nil
}

func (u *PaymentUsecase) GetTransactionByRef(refID string) (*TransferResponse, error) {
	tx, err := u.repo.GetTransactionByRef(refID)
	if err != nil {
		return nil, err
	}

	return &TransferResponse{
		TransactionID: tx.ID,
		Reference:     tx.Reference,
		Amount:        tx.Amount,
		Status:        tx.Status,
		CreatedAt:     tx.CreatedAt,
	}, nil
}

func (u *PaymentUsecase) TopUpWallet(req TopUpRequest) (*TopUpResponse, error) {
	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	tx, err := u.repo.BeginTx()
	if err != nil {
		return nil, err
	}
	defer u.repo.RollbackTx(tx)

	err = u.repo.TopUpWallet(tx, req.UserID, req.Amount)
	if err != nil {
		return nil, err
	}

	wallet, err := u.repo.GetWalletForUpdate(tx, req.UserID) // Get the updated wallet to return the new balance
	if err != nil {
		return nil, err
	}

	err = u.repo.CommitTx(tx)
	if err != nil {
		return nil, err
	}

	return &TopUpResponse{
		UserID:  req.UserID,
		Amount:  req.Amount,
		Balance: wallet.Balance,
	}, nil
}

func (u *PaymentUsecase) GetWallet(userID string) (*GetWalletResponse, error) {
	wallet, err := u.repo.GetWalletByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &GetWalletResponse{
		WalletID:  wallet.ID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		Version:   wallet.Version,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}, nil
}

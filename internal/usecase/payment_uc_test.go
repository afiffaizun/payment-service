
package usecase

import (
	"errors"
	"payment-service/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionRepository is a mock implementation of domain.TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) GetWalletForUpdate(tx interface{}, userID string) (*domain.Wallet, error) {
	args := m.Called(tx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockTransactionRepository) UpdateWalletBalance(tx interface{}, walletID string, amount int64) error {
	args := m.Called(tx, walletID, amount)
	return args.Error(0)
}

func (m *MockTransactionRepository) CreateTransaction(tx interface{}, transaction *domain.Transaction) error {
	args := m.Called(tx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetTransactionByRef(refID string) (*domain.Transaction, error) {
	args := m.Called(refID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) BeginTx() (interface{}, error) {
	args := m.Called()
	return args.Get(0), args.Error(1)
}

func (m *MockTransactionRepository) CommitTx(tx interface{}) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) RollbackTx(tx interface{}) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) TopUpWallet(tx interface{}, userID string, amount int64) error {
	args := m.Called(tx, userID, amount)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetWalletByUserID(userID string) (*domain.Wallet, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func TestTopUpWallet(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	uc := NewPaymentUsecase(mockRepo)

	// Mock for transaction object (can be any non-nil value)
	mockTx := &struct{}{}

	tests := []struct {
		name string
		req  TopUpRequest
		mock func()
		want *TopUpResponse
		err  error
	}{
		{
			name: "Successful TopUp",
			req: TopUpRequest{
				UserID: "111",
				Amount: 1000,
			},
			mock: func() {
				mockRepo.On("BeginTx").Return(mockTx, nil).Once()
				mockRepo.On("TopUpWallet", mockTx, "111", int64(1000)).Return(nil).Once()
				mockRepo.On("GetWalletForUpdate", mockTx, "111").Return(&domain.Wallet{
					ID:      "wallet-111",
					UserID:  "111",
					Balance: 11000,
					Version: 1,
				}, nil).Once()
				mockRepo.On("CommitTx", mockTx).Return(nil).Once()
				mockRepo.On("RollbackTx", mockTx).Return(nil).Once()
			},
			want: &TopUpResponse{
				UserID:  "111",
				Amount:  1000,
				Balance: 11000,
			},
			err: nil,
		},
		{
			name: "Invalid Amount",
			req: TopUpRequest{
				UserID: "111",
				Amount: -100,
			},
			mock: func() {
				// No repository calls expected
			},
			want: nil,
			err:  ErrInvalidAmount,
		},
		{
			name: "BeginTx Error",
			req: TopUpRequest{
				UserID: "111",
				Amount: 1000,
			},
			mock: func() {
				mockRepo.On("BeginTx").Return(nil, errors.New("db error")).Once()
			},
			want: nil,
			err:  errors.New("db error"),
		},
		{
			name: "TopUpWallet Error",
			req: TopUpRequest{
				UserID: "111",
				Amount: 1000,
			},
			mock: func() {
				mockRepo.On("BeginTx").Return(mockTx, nil).Once()
				mockRepo.On("TopUpWallet", mockTx, "111", int64(1000)).Return(errors.New("repo error")).Once()
				mockRepo.On("RollbackTx", mock.Anything).Return(nil).Once()
			},
			want: nil,
			err:  errors.New("repo error"),
		},
		{
			name: "GetWalletForUpdate Error after TopUp",
			req: TopUpRequest{
				UserID: "111",
				Amount: 1000,
			},
			mock: func() {
				mockRepo.On("BeginTx").Return(mockTx, nil).Once()
				mockRepo.On("TopUpWallet", mockTx, "111", int64(1000)).Return(nil).Once()
				mockRepo.On("GetWalletForUpdate", mockTx, "111").Return(nil, errors.New("wallet not found")).Once()
				mockRepo.On("RollbackTx", mock.Anything).Return(nil).Once()
			},
			want: nil,
			err:  errors.New("wallet not found"),
		},
		{
			name: "CommitTx Error",
			req: TopUpRequest{
				UserID: "111",
				Amount: 1000,
			},
			mock: func() {
				mockRepo.On("BeginTx").Return(mockTx, nil).Once()
				mockRepo.On("TopUpWallet", mockTx, "111", int64(1000)).Return(nil).Once()
				mockRepo.On("GetWalletForUpdate", mockTx, "111").Return(&domain.Wallet{
					ID:      "wallet-111",
					UserID:  "111",
					Balance: 11000,
					Version: 1,
				}, nil).Once()
				mockRepo.On("CommitTx", mockTx).Return(errors.New("commit error")).Once()
				mockRepo.On("RollbackTx", mock.Anything).Return(nil).Once()
			},
			want: nil,
			err:  errors.New("commit error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Calls = []mock.Call{} // Clear previous mocks
			tt.mock()

			got, err := uc.TopUpWallet(tt.req)

			if tt.err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.UserID, got.UserID)
				assert.Equal(t, tt.want.Amount, got.Amount)
				assert.Equal(t, tt.want.Balance, got.Balance)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

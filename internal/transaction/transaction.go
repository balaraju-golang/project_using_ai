package transaction

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeDeposit   TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypeTransfer   TransactionType = "transfer"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
	ID            string            `json:"id"`
	Type          TransactionType   `json:"type"`
	AccountID     string            `json:"account_id,omitempty"`
	FromAccountID string            `json:"from_account_id,omitempty"`
	ToAccountID   string            `json:"to_account_id,omitempty"`
	Amount        int64             `json:"amount"`
	Timestamp     time.Time         `json:"timestamp"`
	Status        TransactionStatus `json:"status"`
}

type DepositRequest struct {
	AccountID string `json:"account_id"`
	Amount    int64  `json:"amount"`
}

type WithdrawRequest struct {
	AccountID string `json:"account_id"`
	Amount    int64  `json:"amount"`
}

type TransferRequest struct {
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        int64  `json:"amount"`
}

type TransactionResponse struct {
	TransactionID string            `json:"transaction_id"`
	Status        TransactionStatus `json:"status"`
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateDepositTransaction(accountID string, amount int64) *Transaction {
	return &Transaction{
		ID:        uuid.New().String(),
		Type:      TransactionTypeDeposit,
		AccountID: accountID,
		Amount:    amount,
		Timestamp: time.Now(),
		Status:    TransactionStatusCompleted,
	}
}

func (s *Service) CreateWithdrawalTransaction(accountID string, amount int64) *Transaction {
	return &Transaction{
		ID:        uuid.New().String(),
		Type:      TransactionTypeWithdrawal,
		AccountID: accountID,
		Amount:    amount,
		Timestamp: time.Now(),
		Status:    TransactionStatusCompleted,
	}
}

func (s *Service) CreateTransferTransaction(fromAccountID, toAccountID string, amount int64) *Transaction {
	return &Transaction{
		ID:            uuid.New().String(),
		Type:          TransactionTypeTransfer,
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
		Timestamp:     time.Now(),
		Status:        TransactionStatusCompleted,
	}
}

func (s *Service) CreateFailedTransaction(txType TransactionType, accountID string, amount int64, reason string) *Transaction {
	return &Transaction{
		ID:        uuid.New().String(),
		Type:      txType,
		AccountID: accountID,
		Amount:    amount,
		Timestamp: time.Now(),
		Status:    TransactionStatusFailed,
	}
} 
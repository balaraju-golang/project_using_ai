package account

import (
	"time"

	"banking-service/pkg/errors"

	"github.com/google/uuid"
)

type Account struct {
	ID           string    `json:"id"`
	CustomerName string    `json:"owner_name"`
	Balance      int64     `json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateAccountRequest struct {
	CustomerName   string `json:"customer_name"`
	InitialBalance int64  `json:"initial_balance"`
}

type CreateAccountResponse struct {
	ID           string    `json:"id"`
	CustomerName string    `json:"owner_name"`
	Balance      int64     `json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateAccount(req CreateAccountRequest) (*Account, error) {
	if req.CustomerName == "" {
		return nil, &errors.ErrInvalidCustomerName{Name: req.CustomerName}
	}

	if req.InitialBalance < 0 {
		return nil, &errors.ErrInvalidInitialBalance{Balance: req.InitialBalance}
	}

	accountID := uuid.New().String()
	now := time.Now()
	
	account := &Account{
		ID:           accountID,
		CustomerName: req.CustomerName,
		Balance:      req.InitialBalance,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return account, nil
}

func (s *Service) ValidateAccount(account *Account) error {
	if account == nil {
		return &errors.ErrAccountNotFound{AccountID: "nil"}
	}
	return nil
}

func (s *Service) CanWithdraw(account *Account, amount int64) error {
	if err := s.ValidateAccount(account); err != nil {
		return err
	}

	if amount <= 0 {
		return &errors.ErrInvalidAmount{Amount: amount}
	}

	if account.Balance < amount {
		return &errors.ErrInsufficientFunds{
			AccountID: account.ID,
			Balance:   account.Balance,
			Amount:    amount,
		}
	}

	return nil
}

func (s *Service) Deposit(account *Account, amount int64) error {
	if err := s.ValidateAccount(account); err != nil {
		return err
	}

	if amount <= 0 {
		return &errors.ErrInvalidAmount{Amount: amount}
	}

	account.Balance += amount
	account.UpdatedAt = time.Now()
	return nil
}

func (s *Service) Withdraw(account *Account, amount int64) error {
	if err := s.CanWithdraw(account, amount); err != nil {
		return err
	}

	account.Balance -= amount
	account.UpdatedAt = time.Now()
	return nil
}

func (s *Service) Transfer(fromAccount, toAccount *Account, amount int64) error {
	if err := s.ValidateAccount(fromAccount); err != nil {
		return err
	}
	if err := s.ValidateAccount(toAccount); err != nil {
		return err
	}

	if fromAccount.ID == toAccount.ID {
		return &errors.ErrSameAccountTransfer{
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
		}
	}

	if amount <= 0 {
		return &errors.ErrInvalidAmount{Amount: amount}
	}

	if err := s.CanWithdraw(fromAccount, amount); err != nil {
		return err
	}

	fromAccount.Balance -= amount
	toAccount.Balance += amount
	
	now := time.Now()
	fromAccount.UpdatedAt = now
	toAccount.UpdatedAt = now

	return nil
} 
package errors

import "fmt"

type ErrAccountNotFound struct {
	AccountID string
}

func (e ErrAccountNotFound) Error() string {
	return fmt.Sprintf("account not found: %s", e.AccountID)
}

type ErrInsufficientFunds struct {
	AccountID string
	Balance   int64
	Amount    int64
}

func (e ErrInsufficientFunds) Error() string {
	return fmt.Sprintf("insufficient funds in account %s: balance %d, requested %d", e.AccountID, e.Balance, e.Amount)
}

type ErrInvalidAmount struct {
	Amount int64
}

func (e ErrInvalidAmount) Error() string {
	return fmt.Sprintf("invalid amount: %d (must be positive)", e.Amount)
}

type ErrInvalidCustomerName struct {
	Name string
}

func (e ErrInvalidCustomerName) Error() string {
	return fmt.Sprintf("invalid customer name: %s (must not be empty)", e.Name)
}

type ErrInvalidInitialBalance struct {
	Balance int64
}

func (e ErrInvalidInitialBalance) Error() string {
	return fmt.Sprintf("invalid initial balance: %d (must be non-negative)", e.Balance)
}

type ErrTransactionFailed struct {
	TransactionID string
	Reason        string
}

func (e ErrTransactionFailed) Error() string {
	return fmt.Sprintf("transaction %s failed: %s", e.TransactionID, e.Reason)
}

type ErrSameAccountTransfer struct {
	FromAccountID string
	ToAccountID   string
}

func (e ErrSameAccountTransfer) Error() string {
	return fmt.Sprintf("cannot transfer to same account: from %s to %s", e.FromAccountID, e.ToAccountID)
} 
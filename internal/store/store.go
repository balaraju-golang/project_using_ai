package store

import (
	"fmt"
	"sync"

	"banking-service/internal/account"
	"banking-service/internal/transaction"
	"banking-service/pkg/errors"
)

type Store struct {
	accounts     map[string]*account.Account
	transactions map[string]*transaction.Transaction
	mu           sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		accounts:     make(map[string]*account.Account),
		transactions: make(map[string]*transaction.Transaction),
	}
}

func (s *Store) CreateAccount(acc *account.Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.accounts[acc.ID]; exists {
		return fmt.Errorf("account with ID %s already exists", acc.ID)
	}

	s.accounts[acc.ID] = acc
	return nil
}

func (s *Store) GetAccount(id string) (*account.Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	acc, exists := s.accounts[id]
	if !exists {
		return nil, &errors.ErrAccountNotFound{AccountID: id}
	}

	return acc, nil
}

func (s *Store) UpdateAccount(acc *account.Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.accounts[acc.ID]; !exists {
		return &errors.ErrAccountNotFound{AccountID: acc.ID}
	}

	s.accounts[acc.ID] = acc
	return nil
}

func (s *Store) StoreTransaction(tx *transaction.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.transactions[tx.ID]; exists {
		return fmt.Errorf("transaction with ID %s already exists", tx.ID)
	}

	s.transactions[tx.ID] = tx
	return nil
}

func (s *Store) GetTransaction(id string) (*transaction.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, exists := s.transactions[id]
	if !exists {
		return nil, fmt.Errorf("transaction not found: %s", id)
	}

	return tx, nil
}

func (s *Store) GetAllAccounts() []*account.Account {
	s.mu.RLock()
	defer s.mu.RUnlock()

	accounts := make([]*account.Account, 0, len(s.accounts))
	for _, acc := range s.accounts {
		accounts = append(accounts, acc)
	}
	return accounts
}

func (s *Store) GetAllTransactions() []*transaction.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	transactions := make([]*transaction.Transaction, 0, len(s.transactions))
	for _, tx := range s.transactions {
		transactions = append(transactions, tx)
	}
	return transactions
}

func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.accounts = make(map[string]*account.Account)
	s.transactions = make(map[string]*transaction.Transaction)
} 
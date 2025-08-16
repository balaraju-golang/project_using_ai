package store

import (
	"testing"

	"banking-service/internal/account"
	"banking-service/internal/transaction"
)

func TestNewStore(t *testing.T) {
	store := NewStore()
	if store == nil {
		t.Error("NewStore() returned nil")
	}
}

func TestCreateAccount(t *testing.T) {
	store := NewStore()

	account := &account.Account{
		ID:           "test-id",
		CustomerName: "Ravi Kumar",
		Balance:      1000,
	}

	err := store.CreateAccount(account)
	if err != nil {
		t.Errorf("CreateAccount() error = %v", err)
	}

	retrievedAccount, err := store.GetAccount("test-id")
	if err != nil {
		t.Errorf("GetAccount() error = %v", err)
	}

	if retrievedAccount.ID != account.ID {
		t.Errorf("GetAccount() ID = %v, want %v", retrievedAccount.ID, account.ID)
	}
}

func TestGetAccount(t *testing.T) {
	store := NewStore()

	account := &account.Account{
		ID:           "test-id",
		CustomerName: "Priya",
		Balance:      1000,
	}

	store.CreateAccount(account)

	retrievedAccount, err := store.GetAccount("test-id")
	if err != nil {
		t.Errorf("GetAccount() error = %v", err)
	}

	if retrievedAccount.ID != account.ID {
		t.Errorf("GetAccount() ID = %v, want %v", retrievedAccount.ID, account.ID)
	}

	_, err = store.GetAccount("non-existent")
	if err == nil {
		t.Error("GetAccount() expected error for non-existent account")
	}
}

func TestUpdateAccount(t *testing.T) {
	store := NewStore()

	account := &account.Account{
		ID:           "test-id",
		CustomerName: "Sunil",
		Balance:      1000,
	}

	store.CreateAccount(account)

	account.Balance = 2000
	err := store.UpdateAccount(account)
	if err != nil {
		t.Errorf("UpdateAccount() error = %v", err)
	}

	retrievedAccount, err := store.GetAccount("test-id")
	if err != nil {
		t.Errorf("GetAccount() error = %v", err)
	}

	if retrievedAccount.Balance != 2000 {
		t.Errorf("UpdateAccount() balance = %v, want %v", retrievedAccount.Balance, 2000)
	}
}

func TestStoreTransaction(t *testing.T) {
	store := NewStore()

	transaction := &transaction.Transaction{
		ID:        "test-tx-id",
		Type:      transaction.TransactionTypeDeposit,
		AccountID: "test-account-id",
		Amount:    1000,
	}

	err := store.StoreTransaction(transaction)
	if err != nil {
		t.Errorf("StoreTransaction() error = %v", err)
	}

	retrievedTransaction, err := store.GetTransaction("test-tx-id")
	if err != nil {
		t.Errorf("GetTransaction() error = %v", err)
	}

	if retrievedTransaction.ID != transaction.ID {
		t.Errorf("GetTransaction() ID = %v, want %v", retrievedTransaction.ID, transaction.ID)
	}
}

func TestGetTransaction(t *testing.T) {
	store := NewStore()

	transaction := &transaction.Transaction{
		ID:        "test-tx-id",
		Type:      transaction.TransactionTypeDeposit,
		AccountID: "test-account-id",
		Amount:    1000,
	}

	store.StoreTransaction(transaction)

	retrievedTransaction, err := store.GetTransaction("test-tx-id")
	if err != nil {
		t.Errorf("GetTransaction() error = %v", err)
	}

	if retrievedTransaction.ID != transaction.ID {
		t.Errorf("GetTransaction() ID = %v, want %v", retrievedTransaction.ID, transaction.ID)
	}

	_, err = store.GetTransaction("non-existent")
	if err == nil {
		t.Error("GetTransaction() expected error for non-existent transaction")
	}
}

func TestGetAllAccounts(t *testing.T) {
	store := NewStore()

	account1 := &account.Account{
		ID:           "test-id-1",
		CustomerName: "Ravi Kumar",
		Balance:      1000,
	}

	account2 := &account.Account{
		ID:           "test-id-2",
		CustomerName: "Priya",
		Balance:      2000,
	}

	store.CreateAccount(account1)
	store.CreateAccount(account2)

	accounts := store.GetAllAccounts()
	if len(accounts) != 2 {
		t.Errorf("GetAllAccounts() returned %d accounts, want 2", len(accounts))
	}
}

func TestGetAllTransactions(t *testing.T) {
	store := NewStore()

	transaction1 := &transaction.Transaction{
		ID:        "test-tx-id-1",
		Type:      transaction.TransactionTypeDeposit,
		AccountID: "test-account-id",
		Amount:    1000,
	}

	transaction2 := &transaction.Transaction{
		ID:        "test-tx-id-2",
		Type:      transaction.TransactionTypeWithdrawal,
		AccountID: "test-account-id",
		Amount:    500,
	}

	store.StoreTransaction(transaction1)
	store.StoreTransaction(transaction2)

	transactions := store.GetAllTransactions()
	if len(transactions) != 2 {
		t.Errorf("GetAllTransactions() returned %d transactions, want 2", len(transactions))
	}
}

func TestClear(t *testing.T) {
	store := NewStore()

	account := &account.Account{
		ID:           "test-id",
		CustomerName: "Ravi Kumar",
		Balance:      1000,
	}

	transaction := &transaction.Transaction{
		ID:        "test-tx-id",
		Type:      transaction.TransactionTypeDeposit,
		AccountID: "test-account-id",
		Amount:    1000,
	}

	store.CreateAccount(account)
	store.StoreTransaction(transaction)

	store.Clear()

	accounts := store.GetAllAccounts()
	if len(accounts) != 0 {
		t.Errorf("Clear() accounts not cleared, got %d accounts", len(accounts))
	}

	transactions := store.GetAllTransactions()
	if len(transactions) != 0 {
		t.Errorf("Clear() transactions not cleared, got %d transactions", len(transactions))
	}
}

func TestConcurrency(t *testing.T) {
	store := NewStore()

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			account := &account.Account{
				ID:           "test-id-" + string(rune(id)),
				CustomerName: "Ravi Kumar",
				Balance:      1000,
			}

			err := store.CreateAccount(account)
			if err != nil {
				t.Errorf("CreateAccount() error in goroutine %d: %v", id, err)
			}

			account.Balance = 2000
			err = store.UpdateAccount(account)
			if err != nil {
				t.Errorf("UpdateAccount() error in goroutine %d: %v", id, err)
			}

			_, err = store.GetAccount(account.ID)
			if err != nil {
				t.Errorf("GetAccount() error in goroutine %d: %v", id, err)
			}

			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	accounts := store.GetAllAccounts()
	if len(accounts) != 10 {
		t.Errorf("Concurrency test failed, got %d accounts, want 10", len(accounts))
	}
} 
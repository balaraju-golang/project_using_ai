package account

import (
	"testing"

	"banking-service/pkg/errors"
)

func TestNewService(t *testing.T) {
	service := NewService()
	if service == nil {
		t.Error("NewService() returned nil")
	}
}

func TestCreateAccount(t *testing.T) {
	service := NewService()

	tests := []struct {
		name    string
		req     CreateAccountRequest
		wantErr bool
		errType interface{}
	}{
		{
			name: "valid_account_creation",
			req: CreateAccountRequest{
				CustomerName:   "Ravi Kumar",
				InitialBalance: 10000,
			},
			wantErr: false,
		},
		{
			name: "empty_customer_name",
			req: CreateAccountRequest{
				CustomerName:   "",
				InitialBalance: 10000,
			},
			wantErr:  true,
			errType: &errors.ErrInvalidCustomerName{},
		},
		{
			name: "negative_initial_balance",
			req: CreateAccountRequest{
				CustomerName:   "Priya",
				InitialBalance: -1000,
			},
			wantErr:  true,
			errType: &errors.ErrInvalidInitialBalance{},
		},
		{
			name: "zero_initial_balance",
			req: CreateAccountRequest{
				CustomerName:   "Sunil",
				InitialBalance: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := service.CreateAccount(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Error("CreateAccount() expected error but got none")
					return
				}

				switch tt.errType.(type) {
				case *errors.ErrInvalidCustomerName:
					if _, ok := err.(*errors.ErrInvalidCustomerName); !ok {
						t.Errorf("CreateAccount() error type = %T, want *errors.ErrInvalidCustomerName", err)
					}
				case *errors.ErrInvalidInitialBalance:
					if _, ok := err.(*errors.ErrInvalidInitialBalance); !ok {
						t.Errorf("CreateAccount() error type = %T, want *errors.ErrInvalidInitialBalance", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("CreateAccount() unexpected error = %v", err)
					return
				}

				if account == nil {
					t.Error("CreateAccount() returned nil account")
					return
				}

				if account.CustomerName != tt.req.CustomerName {
					t.Errorf("CreateAccount() customer name = %v, want %v", account.CustomerName, tt.req.CustomerName)
				}

				if account.Balance != tt.req.InitialBalance {
					t.Errorf("CreateAccount() balance = %v, want %v", account.Balance, tt.req.InitialBalance)
				}

				if account.ID == "" {
					t.Error("CreateAccount() returned empty ID")
				}

				if account.CreatedAt.IsZero() {
					t.Error("CreateAccount() returned zero CreatedAt")
				}

				if account.UpdatedAt.IsZero() {
					t.Error("CreateAccount() returned zero UpdatedAt")
				}
			}
		})
	}
}

func TestValidateAccount(t *testing.T) {
	service := NewService()

	tests := []struct {
		name    string
		account *Account
		wantErr bool
	}{
		{
			name: "valid_account",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Anjali",
				Balance:      1000,
			},
			wantErr: false,
		},
		{
			name:    "nil_account",
			account: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateAccount(tt.account)

			if tt.wantErr {
				if err == nil {
					t.Error("ValidateAccount() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateAccount() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestCanWithdraw(t *testing.T) {
	service := NewService()

	tests := []struct {
		name    string
		account *Account
		amount  int64
		wantErr bool
		errType interface{}
	}{
		{
			name: "valid_withdrawal",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Deepak",
				Balance:      1000,
			},
			amount:  500,
			wantErr: false,
		},
		{
			name: "exact_balance_withdrawal",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Suresh",
				Balance:      1000,
			},
			amount:  1000,
			wantErr: false,
		},
		{
			name: "insufficient_funds",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Meena",
				Balance:      1000,
			},
			amount:  1500,
			wantErr: true,
			errType: &errors.ErrInsufficientFunds{},
		},
		{
			name: "negative_amount",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Ramesh",
				Balance:      1000,
			},
			amount:  -100,
			wantErr: true,
			errType: &errors.ErrInvalidAmount{},
		},
		{
			name: "zero_amount",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Neha",
				Balance:      1000,
			},
			amount:  0,
			wantErr: true,
			errType: &errors.ErrInvalidAmount{},
		},
		{
			name:    "nil_account",
			account: nil,
			amount:  100,
			wantErr: true,
			errType: &errors.ErrAccountNotFound{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CanWithdraw(tt.account, tt.amount)

			if tt.wantErr {
				if err == nil {
					t.Error("CanWithdraw() expected error but got none")
					return
				}

				switch tt.errType.(type) {
				case *errors.ErrInsufficientFunds:
					if _, ok := err.(*errors.ErrInsufficientFunds); !ok {
						t.Errorf("CanWithdraw() error type = %T, want *errors.ErrInsufficientFunds", err)
					}
				case *errors.ErrInvalidAmount:
					if _, ok := err.(*errors.ErrInvalidAmount); !ok {
						t.Errorf("CanWithdraw() error type = %T, want *errors.ErrInvalidAmount", err)
					}
				case *errors.ErrAccountNotFound:
					if _, ok := err.(*errors.ErrAccountNotFound); !ok {
						t.Errorf("CanWithdraw() error type = %T, want *errors.ErrAccountNotFound", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("CanWithdraw() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestDeposit(t *testing.T) {
	service := NewService()

	tests := []struct {
		name    string
		account *Account
		amount  int64
		wantErr bool
	}{
		{
			name: "valid_deposit",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Amit",
				Balance:      1000,
			},
			amount:  500,
			wantErr: false,
		},
		{
			name: "zero_deposit",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Kiran",
				Balance:      1000,
			},
			amount:  0,
			wantErr: true,
		},
		{
			name: "negative_deposit",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Pooja",
				Balance:      1000,
			},
			amount:  -100,
			wantErr: true,
		},
		{
			name:    "nil_account",
			account: nil,
			amount:  100,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var originalBalance int64
			if tt.account != nil {
				originalBalance = tt.account.Balance
			}
			err := service.Deposit(tt.account, tt.amount)

			if tt.wantErr {
				if err == nil {
					t.Error("Deposit() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Deposit() unexpected error = %v", err)
					return
				}

				expectedBalance := originalBalance + tt.amount
				if tt.account.Balance != expectedBalance {
					t.Errorf("Deposit() balance = %v, want %v", tt.account.Balance, expectedBalance)
				}

				if tt.account.UpdatedAt.IsZero() {
					t.Error("Deposit() did not update UpdatedAt")
				}
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	service := NewService()

	tests := []struct {
		name    string
		account *Account
		amount  int64
		wantErr bool
	}{
		{
			name: "valid_withdrawal",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Vijay",
				Balance:      1000,
			},
			amount:  500,
			wantErr: false,
		},
		{
			name: "exact_balance_withdrawal",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Arun",
				Balance:      1000,
			},
			amount:  1000,
			wantErr: false,
		},
		{
			name: "insufficient_funds",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Suman",
				Balance:      1000,
			},
			amount:  1500,
			wantErr: true,
		},
		{
			name: "zero_withdrawal",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Lakshmi",
				Balance:      1000,
			},
			amount:  0,
			wantErr: true,
		},
		{
			name: "negative_withdrawal",
			account: &Account{
				ID:           "test-id",
				CustomerName: "Manoj",
				Balance:      1000,
			},
			amount:  -100,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalBalance := tt.account.Balance
			err := service.Withdraw(tt.account, tt.amount)

			if tt.wantErr {
				if err == nil {
					t.Error("Withdraw() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Withdraw() unexpected error = %v", err)
					return
				}

				expectedBalance := originalBalance - tt.amount
				if tt.account.Balance != expectedBalance {
					t.Errorf("Withdraw() balance = %v, want %v", tt.account.Balance, expectedBalance)
				}

				if tt.account.UpdatedAt.IsZero() {
					t.Error("Withdraw() did not update UpdatedAt")
				}
			}
		})
	}
}

func TestTransfer(t *testing.T) {
	service := NewService()

	tests := []struct {
		name         string
		fromAccount  *Account
		toAccount    *Account
		amount       int64
		wantErr      bool
	}{
		{
			name: "valid_transfer",
			fromAccount: &Account{
				ID:           "from-id",
				CustomerName: "Nisha",
				Balance:      1000,
			},
			toAccount: &Account{
				ID:           "to-id",
				CustomerName: "Rahul",
				Balance:      500,
			},
			amount:  300,
			wantErr: false,
		},
		{
			name: "insufficient_funds",
			fromAccount: &Account{
				ID:           "from-id",
				CustomerName: "Shyam",
				Balance:      1000,
			},
			toAccount: &Account{
				ID:           "to-id",
				CustomerName: "Geeta",
				Balance:      500,
			},
			amount:  1500,
			wantErr: true,
		},
		{
			name: "same_account_transfer",
			fromAccount: &Account{
				ID:           "same-id",
				CustomerName: "Ajay",
				Balance:      1000,
			},
			toAccount: &Account{
				ID:           "same-id",
				CustomerName: "Ajay",
				Balance:      1000,
			},
			amount:  300,
			wantErr: true,
		},
		{
			name: "zero_amount_transfer",
			fromAccount: &Account{
				ID:           "from-id",
				CustomerName: "Sneha",
				Balance:      1000,
			},
			toAccount: &Account{
				ID:           "to-id",
				CustomerName: "Tarun",
				Balance:      500,
			},
			amount:  0,
			wantErr: true,
		},
		{
			name: "negative_amount_transfer",
			fromAccount: &Account{
				ID:           "from-id",
				CustomerName: "Manoj",
				Balance:      1000,
			},
			toAccount: &Account{
				ID:           "to-id",
				CustomerName: "Jane Doe",
				Balance:      500,
			},
			amount:  -100,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fromOriginalBalance := tt.fromAccount.Balance
			toOriginalBalance := tt.toAccount.Balance

			err := service.Transfer(tt.fromAccount, tt.toAccount, tt.amount)

			if tt.wantErr {
				if err == nil {
					t.Error("Transfer() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Transfer() unexpected error = %v", err)
					return
				}

				expectedFromBalance := fromOriginalBalance - tt.amount
				if tt.fromAccount.Balance != expectedFromBalance {
					t.Errorf("Transfer() from account balance = %v, want %v", tt.fromAccount.Balance, expectedFromBalance)
				}

				expectedToBalance := toOriginalBalance + tt.amount
				if tt.toAccount.Balance != expectedToBalance {
					t.Errorf("Transfer() to account balance = %v, want %v", tt.toAccount.Balance, expectedToBalance)
				}

				if tt.fromAccount.UpdatedAt.IsZero() {
					t.Error("Transfer() did not update from account UpdatedAt")
				}

				if tt.toAccount.UpdatedAt.IsZero() {
					t.Error("Transfer() did not update to account UpdatedAt")
				}
			}
		})
	}
} 
package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"banking-service/internal/account"
	"banking-service/internal/store"
	"banking-service/internal/transaction"
	"banking-service/pkg/errors"
)

type Handler struct {
	store           *store.Store
	accountService  *account.Service
	transactionService *transaction.Service
	logger          *logrus.Logger
}

func NewHandler(store *store.Store, logger *logrus.Logger) *Handler {
	return &Handler{
		store:             store,
		accountService:    account.NewService(),
		transactionService: transaction.NewService(),
		logger:            logger,
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func (h *Handler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.WithError(err).Error("Failed to encode JSON response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) writeError(w http.ResponseWriter, statusCode int, message string) {
	h.writeJSON(w, statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	var req account.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode create account request")
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	acc, err := h.accountService.CreateAccount(req)
	if err != nil {
		h.logger.WithError(err).WithField("customer_name", req.CustomerName).Error("Failed to create account")
		
		switch err.(type) {
		case *errors.ErrInvalidCustomerName, *errors.ErrInvalidInitialBalance:
			h.writeError(w, http.StatusBadRequest, err.Error())
		default:
			h.writeError(w, http.StatusInternalServerError, "Failed to create account")
		}
		return
	}
	
	if err := h.store.CreateAccount(acc); err != nil {
		h.logger.WithError(err).WithField("account_id", acc.ID).Error("Failed to store account")
		h.writeError(w, http.StatusInternalServerError, "Failed to create account")
		return
	}
	
	h.logger.WithFields(logrus.Fields{
		"account_id": acc.ID,
		"customer_name": acc.CustomerName,
		"balance": acc.Balance,
	}).Info("Account created successfully")
	
	h.writeJSON(w, http.StatusCreated, acc)
}

func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	path := strings.TrimPrefix(r.URL.Path, "/accounts/")
	if path == "" {
		h.writeError(w, http.StatusBadRequest, "Account ID required")
		return
	}
	
	acc, err := h.store.GetAccount(path)
	if err != nil {
		h.logger.WithError(err).WithField("account_id", path).Error("Failed to get account")
		
		if _, ok := err.(*errors.ErrAccountNotFound); ok {
			h.writeError(w, http.StatusNotFound, "Account not found")
		} else {
			h.writeError(w, http.StatusInternalServerError, "Failed to get account")
		}
		return
	}
	
	h.logger.WithField("account_id", path).Info("Account retrieved successfully")
	h.writeJSON(w, http.StatusOK, acc)
}

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	var req transaction.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode deposit request")
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	acc, err := h.store.GetAccount(req.AccountID)
	if err != nil {
		h.logger.WithError(err).WithField("account_id", req.AccountID).Error("Failed to get account for deposit")
		
		if _, ok := err.(*errors.ErrAccountNotFound); ok {
			h.writeError(w, http.StatusNotFound, "Account not found")
		} else {
			h.writeError(w, http.StatusInternalServerError, "Failed to process deposit")
		}
		return
	}
	
	if err := h.accountService.Deposit(acc, req.Amount); err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"account_id": req.AccountID,
			"amount": req.Amount,
		}).Error("Failed to process deposit")
		
		if _, ok := err.(*errors.ErrInvalidAmount); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
		} else {
			h.writeError(w, http.StatusInternalServerError, "Failed to process deposit")
		}
		return
	}
	
	if err := h.store.UpdateAccount(acc); err != nil {
		h.logger.WithError(err).WithField("account_id", acc.ID).Error("Failed to update account after deposit")
		h.writeError(w, http.StatusInternalServerError, "Failed to process deposit")
		return
	}
	
	tx := h.transactionService.CreateDepositTransaction(req.AccountID, req.Amount)
	if err := h.store.StoreTransaction(tx); err != nil {
		h.logger.WithError(err).WithField("transaction_id", tx.ID).Error("Failed to store deposit transaction")
	}
	
	h.logger.WithFields(logrus.Fields{
		"account_id": req.AccountID,
		"amount": req.Amount,
		"transaction_id": tx.ID,
		"new_balance": acc.Balance,
	}).Info("Deposit processed successfully")
	
	h.writeJSON(w, http.StatusOK, transaction.TransactionResponse{
		TransactionID: tx.ID,
		Status:        tx.Status,
	})
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	var req transaction.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode withdraw request")
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	acc, err := h.store.GetAccount(req.AccountID)
	if err != nil {
		h.logger.WithError(err).WithField("account_id", req.AccountID).Error("Failed to get account for withdrawal")
		
		if _, ok := err.(*errors.ErrAccountNotFound); ok {
			h.writeError(w, http.StatusNotFound, "Account not found")
		} else {
			h.writeError(w, http.StatusInternalServerError, "Failed to process withdrawal")
		}
		return
	}
	
	if err := h.accountService.Withdraw(acc, req.Amount); err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"account_id": req.AccountID,
			"amount": req.Amount,
		}).Error("Failed to process withdrawal")
		
		switch err.(type) {
		case *errors.ErrInvalidAmount:
			h.writeError(w, http.StatusBadRequest, err.Error())
		case *errors.ErrInsufficientFunds:
			h.writeError(w, http.StatusBadRequest, err.Error())
		default:
			h.writeError(w, http.StatusInternalServerError, "Failed to process withdrawal")
		}
		return
	}
	
	if err := h.store.UpdateAccount(acc); err != nil {
		h.logger.WithError(err).WithField("account_id", acc.ID).Error("Failed to update account after withdrawal")
		h.writeError(w, http.StatusInternalServerError, "Failed to process withdrawal")
		return
	}
	
	tx := h.transactionService.CreateWithdrawalTransaction(req.AccountID, req.Amount)
	if err := h.store.StoreTransaction(tx); err != nil {
		h.logger.WithError(err).WithField("transaction_id", tx.ID).Error("Failed to store withdrawal transaction")
	}
	
	h.logger.WithFields(logrus.Fields{
		"account_id": req.AccountID,
		"amount": req.Amount,
		"transaction_id": tx.ID,
		"new_balance": acc.Balance,
	}).Info("Withdrawal processed successfully")
	
	h.writeJSON(w, http.StatusOK, transaction.TransactionResponse{
		TransactionID: tx.ID,
		Status:        tx.Status,
	})
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	var req transaction.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode transfer request")
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	fromAccount, err := h.store.GetAccount(req.FromAccountID)
	if err != nil {
		h.logger.WithError(err).WithField("from_account_id", req.FromAccountID).Error("Failed to get from account for transfer")
		
		if _, ok := err.(*errors.ErrAccountNotFound); ok {
			h.writeError(w, http.StatusNotFound, "From account not found")
		} else {
			h.writeError(w, http.StatusInternalServerError, "Failed to process transfer")
		}
		return
	}
	
	toAccount, err := h.store.GetAccount(req.ToAccountID)
	if err != nil {
		h.logger.WithError(err).WithField("to_account_id", req.ToAccountID).Error("Failed to get to account for transfer")
		
		if _, ok := err.(*errors.ErrAccountNotFound); ok {
			h.writeError(w, http.StatusNotFound, "To account not found")
		} else {
			h.writeError(w, http.StatusInternalServerError, "Failed to process transfer")
		}
		return
	}
	
	if err := h.accountService.Transfer(fromAccount, toAccount, req.Amount); err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"from_account_id": req.FromAccountID,
			"to_account_id": req.ToAccountID,
			"amount": req.Amount,
		}).Error("Failed to process transfer")
		
		switch err.(type) {
		case *errors.ErrInvalidAmount:
			h.writeError(w, http.StatusBadRequest, err.Error())
		case *errors.ErrInsufficientFunds:
			h.writeError(w, http.StatusBadRequest, err.Error())
		case *errors.ErrSameAccountTransfer:
			h.writeError(w, http.StatusBadRequest, err.Error())
		default:
			h.writeError(w, http.StatusInternalServerError, "Failed to process transfer")
		}
		return
	}
	
	if err := h.store.UpdateAccount(fromAccount); err != nil {
		h.logger.WithError(err).WithField("account_id", fromAccount.ID).Error("Failed to update from account after transfer")
		h.writeError(w, http.StatusInternalServerError, "Failed to process transfer")
		return
	}
	
	if err := h.store.UpdateAccount(toAccount); err != nil {
		h.logger.WithError(err).WithField("account_id", toAccount.ID).Error("Failed to update to account after transfer")
		h.writeError(w, http.StatusInternalServerError, "Failed to process transfer")
		return
	}
	
	tx := h.transactionService.CreateTransferTransaction(req.FromAccountID, req.ToAccountID, req.Amount)
	if err := h.store.StoreTransaction(tx); err != nil {
		h.logger.WithError(err).WithField("transaction_id", tx.ID).Error("Failed to store transfer transaction")
	}
	
	h.logger.WithFields(logrus.Fields{
		"from_account_id": req.FromAccountID,
		"to_account_id": req.ToAccountID,
		"amount": req.Amount,
		"transaction_id": tx.ID,
		"from_balance": fromAccount.Balance,
		"to_balance": toAccount.Balance,
	}).Info("Transfer processed successfully")
	
	h.writeJSON(w, http.StatusOK, transaction.TransactionResponse{
		TransactionID: tx.ID,
		Status:        tx.Status,
	})
} 
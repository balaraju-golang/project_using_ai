package api

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"banking-service/internal/store"
)

type Server struct {
	server *http.Server
	logger *logrus.Logger
	store  *store.Store
}

func NewServer(port string, logger *logrus.Logger, store *store.Store) *Server {
	mux := http.NewServeMux()
	
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	return &Server{
		server: server,
		logger: logger,
		store:  store,
	}
}

func (s *Server) SetupRoutes() {
	handler := NewHandler(s.store, s.logger)
	
	mux := s.server.Handler.(*http.ServeMux)
	mux.HandleFunc("/accounts", handler.CreateAccount)
	mux.HandleFunc("/accounts/", s.getAccountHandler)
	
	mux.HandleFunc("/transactions/deposit", handler.Deposit)
	mux.HandleFunc("/transactions/withdraw", handler.Withdraw)
	mux.HandleFunc("/transactions/transfer", handler.Transfer)
	
	mux.HandleFunc("/health", s.healthCheck)
}

func (s *Server) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	handler := NewHandler(s.store, s.logger)
	handler.GetAccount(w, r)
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func (s *Server) Start() error {
	s.SetupRoutes()
	s.logger.Info("Server starting on port " + s.server.Addr)
	return s.server.ListenAndServe()
} 
package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"banking-service/internal/api"
	"banking-service/internal/store"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	logger.Info("Starting banking service on port " + port)
	
	store := store.NewStore()
	server := api.NewServer(port, logger, store)
	
	if err := server.Start(); err != nil {
		logger.Fatal("Server failed to start: " + err.Error())
	}
} 
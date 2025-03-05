package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jazwu/tiktok-ecommerce/app/user/payment-api/internal/api"
	"github.com/jazwu/tiktok-ecommerce/app/user/payment-api/internal/storage"
)

func main() {
	ctx := context.Background()
	
	// Get database config
	dbConfig := api.GetDatabaseConfig()
	dsn := dbConfig.DSN()

	// Initialize MySQL storage
	dbStorage, err := storage.NewMySQLStorage(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbStorage.Close()

	// Create router with dependencies
	router := api.NewRouter(dbStorage)

	// Configure server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("Starting payment API server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

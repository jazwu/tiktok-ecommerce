package api

import (
	"github.com/gorilla/mux"
	"github.com/jazwu/tiktok-ecommerce/app/user/payment-api/internal/api/handlers"
	"github.com/jazwu/tiktok-ecommerce/app/user/payment-api/internal/api/middleware"
	"github.com/jazwu/tiktok-ecommerce/app/user/payment-api/internal/storage"
)

func NewRouter(storage storage.MySQLStorage) *mux.Router {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.ContentTypeJSONMiddleware)

	// Initialize handlers
	paymentHandler := handlers.NewPaymentHandler(storage)

	// Payment API routes
	router.HandleFunc("/api/v1/payments", paymentHandler.CreatePayment).Methods("POST")
	router.HandleFunc("/api/v1/payments/{id}", paymentHandler.GetPayment).Methods("GET")

	return router
}

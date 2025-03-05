package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jazwu/tiktok-ecommerce/app/user/payment-api/internal/models"
	"github.com/jazwu/tiktok-ecommerce/app/user/payment-api/internal/storage"
)

type PaymentHandler struct {
	storage storage.MySQLStorage
}

func NewPaymentHandler(storage storage.MySQLStorage) *PaymentHandler {
	return &PaymentHandler{storage: storage}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request body
	var req models.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the request
	if req.Amount <= 0 || req.Currency == "" || req.PaymentMethod == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid payment request")
		return
	}

	// Check and store idempotency key
	if err := h.storage.CheckAndStoreIdempotencyKey(ctx, req.IdempotencyKey); err != nil {
		if errors.Is(err, storage.ErrIdempotencyConflict) {
			respondWithError(w, http.StatusConflict, "Duplicate idempotency key")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Create a new payment intent
	intent := models.NewPaymentIntent(req.Amount, req.Currency, req.PaymentMethod)

	// Process the payment (simulated)
	if err := processPayment(intent); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Store the payment intent in the database
	if err := h.storage.StorePaymentIntent(ctx, intent); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to store payment")
		return
	}

	// Return the created payment intent
	respondWithJSON(w, http.StatusCreated, intent)
}

func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paymentID := mux.Vars(r)["id"]

	// Retrieve the payment intent from the database
	intent, err := h.storage.GetPaymentIntent(ctx, paymentID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Payment not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Return the payment intent
	respondWithJSON(w, http.StatusOK, intent)
}

// processPayment simulates payment processing logic
func processPayment(intent *models.PaymentIntent) error {
	switch intent.PaymentMethod {
	case "card":
		if intent.Amount > 10000 { // Simulated fraud check
			intent.Status = models.StatusFailed
			return errors.New("payment declined")
		}
		intent.Status = models.StatusSucceeded
		return nil
	default:
		intent.Status = models.StatusFailed
		return errors.New("unsupported payment method")
	}
}

// respondWithJSON is a helper function to send JSON responses
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError is a helper function to send error responses
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, models.APIError{Error: message})
}

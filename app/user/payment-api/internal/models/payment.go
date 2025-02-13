package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentIntent struct {
	ID            string    `json:"id"`
	Amount        int       `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
}

type PaymentRequest struct {
	Amount         int    `json:"amount" validate:"required,gt=0"`
	Currency       string `json:"currency" validate:"required,len=3"`
	PaymentMethod  string `json:"payment_method" validate:"required"`
	IdempotencyKey string `json:"idempotency_key" validate:"required"`
}

type APIError struct {
	Error string `json:"error"`
}

const (
	StatusRequiresAction = "requires_action"
	StatusSucceeded      = "succeeded"
	StatusFailed         = "failed"
)

func NewPaymentIntent(amount int, currency, paymentMethod string) *PaymentIntent {
	return &PaymentIntent{
		ID:            uuid.New().String(),
		Amount:        amount,
		Currency:      currency,
		PaymentMethod: paymentMethod,
		CreatedAt:     time.Now().UTC(),
		Status:        StatusRequiresAction,
	}
}

package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jazwu/tiktok-ecommerce/app/user/payment-api/internal/models"
)

type MySQLStorage struct {
	db *sql.DB
}

var (
	ErrNotFound            = errors.New("not found")
	ErrIdempotencyConflict = errors.New("idempotency key conflict")
)

// Initialize and return a new instance of MySQLStroge
func NewMySQLStorage(ctx context.Context, dsn string) (*MySQLStorage, error) {
	// database pool initialization and configuration
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open databases: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &MySQLStorage{db: db}, nil
}

// Store a payment intent record in the database using a transaction
func (s *MySQLStorage) StorePaymentIntent(ctx context.Context, intent *models.PaymentIntent) error {
	// Begin a database transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Ensure rollback is called in case of failure
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO payment_intents
	(id, amount, currency, status, payment_method, created_at)
	VALUES (?, ?, ?, ?, ?, ?)`,
		intent.ID, intent.Amount, intent.Currency,
		intent.Status, intent.PaymentMethod, intent.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("errot storing payment intent: %w", err)
	}

	// Commit the transaction to save changes
	return tx.Commit()
}

// Retrieve a payment intent from the database by its ID
func (s *MySQLStorage) GetPaymentIntent(ctx context.Context, id string) (*models.PaymentIntent, error) {
	var intent models.PaymentIntent

	// Execute the SQL query to fetch the payment intent by ID
	err := s.db.QueryRowContext(ctx,
		`SELECT id, amount, currency, status, payment_method, created_at 
	FROM payment_intents WHERE id = ?`,
		id,
	).Scan(
		&intent.ID,
		&intent.Amount,
		&intent.Currency,
		&intent.Status,
		&intent.PaymentMethod,
		&intent.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching payment intent: %w", err)
	}

	return &intent, nil
}

// Ensure an operation is idempotent by attempting to store an idempotency key in the database
func (s *MySQLStorage) CheckAndStoreIdempotencyKey(ctx context.Context, key string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Try to insert the key
	res, err := tx.ExecContext(ctx,
		`INSERT INTO idempotency_keys (idempotency_key) VALUES (?)`,
		key,
	)

	if err != nil {
		if isDuplicationKeyError(err) {
			return ErrIdempotencyConflict
		}
		return fmt.Errorf("error storing idempotency key: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrIdempotencyConflict
	}

	return tx.Commit()
}

// Detect MySQL error code 1062 (ER_DUP_ENTRY)
func isDuplicationKeyError(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return mysqlErr.Number == 1062
	}
	return false
}

func (s *MySQLStorage) Close() error {
	return s.db.Close()
}

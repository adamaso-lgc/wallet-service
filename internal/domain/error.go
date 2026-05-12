package domain

import "errors"

var (
	ErrOwnerIDRequired        = errors.New("ownerID is required")
	ErrInvalidAmount          = errors.New("amount must be greater than zero")
	ErrCurrencyRequired       = errors.New("currency is required")
	ErrCurrencyMismatch       = errors.New("currency mismatch")
	ErrInsufficientFunds      = errors.New("insufficient funds")
	ErrWalletNotActive        = errors.New("wallet is not active")
	ErrUnsupportedEvent       = errors.New("unsupported event type")
	ErrNotFound               = errors.New("aggregate not found")
	ErrConcurrentModification = errors.New("concurrent modification detected")
)

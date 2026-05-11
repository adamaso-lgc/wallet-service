package domain_test

import (
	"errors"
	"testing"

	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		currency string
		wantErr  error
	}{
		{name: "valid money", amount: 100, currency: "USD"},
		{name: "zero amount", amount: 0, currency: "USD", wantErr: domain.ErrInvalidAmount},
		{name: "negative amount", amount: -10, currency: "USD", wantErr: domain.ErrInvalidAmount},
		{name: "empty currency", amount: 50, currency: "", wantErr: domain.ErrCurrencyRequired},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := domain.NewMoney(tc.amount, tc.currency)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.amount, m.Amount())
			assert.Equal(t, tc.currency, m.Currency())
		})
	}
}

func TestMoneyAdd(t *testing.T) {
	tests := []struct {
		name       string
		base       domain.Money
		other      domain.Money
		wantAmount float64
		wantErr    error
	}{
		{
			name:       "same currency",
			base:       mustMoney(t, 100, "USD"),
			other:      mustMoney(t, 50, "USD"),
			wantAmount: 150,
		},
		{
			name:    "currency mismatch",
			base:    mustMoney(t, 100, "USD"),
			other:   mustMoney(t, 50, "EUR"),
			wantErr: domain.ErrCurrencyMismatch,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.base.Add(tc.other)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantAmount, result.Amount())
		})
	}
}

func TestMoneySubtract(t *testing.T) {
	tests := []struct {
		name       string
		base       domain.Money
		other      domain.Money
		wantAmount float64
		wantErr    error
	}{
		{
			name:       "sufficient funds",
			base:       mustMoney(t, 100, "USD"),
			other:      mustMoney(t, 40, "USD"),
			wantAmount: 60,
		},
		{
			name:    "insufficient funds",
			base:    mustMoney(t, 30, "USD"),
			other:   mustMoney(t, 50, "USD"),
			wantErr: domain.ErrInsufficientFunds,
		},
		{
			name:    "currency mismatch",
			base:    mustMoney(t, 100, "USD"),
			other:   mustMoney(t, 50, "EUR"),
			wantErr: domain.ErrCurrencyMismatch,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.base.Subtract(tc.other)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantAmount, result.Amount())
		})
	}
}

// mustMoney is a test helper that creates a Money value or fails the test immediately.
func mustMoney(t *testing.T, amount float64, currency string) domain.Money {
	t.Helper()
	m, err := domain.NewMoney(amount, currency)
	require.NoError(t, err)
	return m
}

// Compile-time check to ensure the errors package is used (errors.Is).
var _ = errors.Is

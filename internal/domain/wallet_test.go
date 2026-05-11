package domain_test

import (
	"testing"

	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWallet(t *testing.T) {
	tests := []struct {
		name           string
		ownerID        string
		currency       string
		initialBalance float64
		wantErr        error
	}{
		{name: "valid wallet", ownerID: "owner-1", currency: "USD", initialBalance: 100},
		{name: "empty ownerID", ownerID: "", currency: "USD", initialBalance: 100, wantErr: domain.ErrOwnerIDRequired},
		{name: "invalid amount", ownerID: "owner-1", currency: "USD", initialBalance: -5, wantErr: domain.ErrInvalidAmount},
		{name: "empty currency", ownerID: "owner-1", currency: "", initialBalance: 100, wantErr: domain.ErrCurrencyRequired},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w, err := domain.NewWallet(tc.ownerID, tc.currency, tc.initialBalance)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, w.GetID())
			assert.Equal(t, int64(1), w.GetVersion())
		})
	}
}

func TestWalletDeposit(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T) *domain.Wallet
		amount    float64
		reference string
		wantErr   error
	}{
		{
			name:      "successful deposit",
			setup:     func(t *testing.T) *domain.Wallet { return mustWallet(t, "owner-1", "USD", 100) },
			amount:    50,
			reference: "ref-1",
		},
		{
			name:      "invalid deposit amount",
			setup:     func(t *testing.T) *domain.Wallet { return mustWallet(t, "owner-1", "USD", 100) },
			amount:    -10,
			reference: "ref-2",
			wantErr:   domain.ErrInvalidAmount,
		},
		{
			name: "deposit on frozen wallet",
			setup: func(t *testing.T) *domain.Wallet {
				w := mustWallet(t, "owner-1", "USD", 100)
				require.NoError(t, w.Freeze("freeze-ref"))
				return w
			},
			amount:    50,
			reference: "ref-3",
			wantErr:   domain.ErrWalletNotActive,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := tc.setup(t)
			err := w.Deposit(tc.amount, tc.reference)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestWalletWithdraw(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T) *domain.Wallet
		amount    float64
		reference string
		wantErr   error
	}{
		{
			name:      "successful withdrawal",
			setup:     func(t *testing.T) *domain.Wallet { return mustWallet(t, "owner-1", "USD", 100) },
			amount:    40,
			reference: "ref-1",
		},
		{
			name:      "insufficient funds",
			setup:     func(t *testing.T) *domain.Wallet { return mustWallet(t, "owner-1", "USD", 30) },
			amount:    50,
			reference: "ref-2",
			wantErr:   domain.ErrInsufficientFunds,
		},
		{
			name:      "invalid amount",
			setup:     func(t *testing.T) *domain.Wallet { return mustWallet(t, "owner-1", "USD", 100) },
			amount:    -5,
			reference: "ref-3",
			wantErr:   domain.ErrInvalidAmount,
		},
		{
			name: "withdraw from frozen wallet",
			setup: func(t *testing.T) *domain.Wallet {
				w := mustWallet(t, "owner-1", "USD", 100)
				require.NoError(t, w.Freeze("freeze-ref"))
				return w
			},
			amount:    10,
			reference: "ref-4",
			wantErr:   domain.ErrWalletNotActive,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := tc.setup(t)
			err := w.Withdraw(tc.amount, tc.reference)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestWalletDebitForTransfer(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T) *domain.Wallet
		amount         float64
		counterpartyID string
		reference      string
		wantErr        error
	}{
		{
			name:           "successful debit",
			setup:          func(t *testing.T) *domain.Wallet { return mustWallet(t, "owner-1", "USD", 100) },
			amount:         60,
			counterpartyID: "owner-2",
			reference:      "ref-1",
		},
		{
			name:           "insufficient funds",
			setup:          func(t *testing.T) *domain.Wallet { return mustWallet(t, "owner-1", "USD", 30) },
			amount:         50,
			counterpartyID: "owner-2",
			reference:      "ref-2",
			wantErr:        domain.ErrInsufficientFunds,
		},
		{
			name: "debit on frozen wallet",
			setup: func(t *testing.T) *domain.Wallet {
				w := mustWallet(t, "owner-1", "USD", 100)
				require.NoError(t, w.Freeze("freeze-ref"))
				return w
			},
			amount:         10,
			counterpartyID: "owner-2",
			reference:      "ref-3",
			wantErr:        domain.ErrWalletNotActive,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := tc.setup(t)
			err := w.DebitForTransfer(tc.amount, tc.counterpartyID, tc.reference)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestWalletCreditForTransfer(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T) *domain.Wallet
		amount         float64
		counterpartyID string
		reference      string
		wantErr        error
	}{
		{
			name:           "successful credit",
			setup:          func(t *testing.T) *domain.Wallet { return mustWallet(t, "owner-1", "USD", 50) },
			amount:         25,
			counterpartyID: "owner-2",
			reference:      "ref-1",
		},
		{
			name: "credit on frozen wallet",
			setup: func(t *testing.T) *domain.Wallet {
				w := mustWallet(t, "owner-1", "USD", 100)
				require.NoError(t, w.Freeze("freeze-ref"))
				return w
			},
			amount:         10,
			counterpartyID: "owner-2",
			reference:      "ref-2",
			wantErr:        domain.ErrWalletNotActive,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := tc.setup(t)
			err := w.CreditForTransfer(tc.amount, tc.counterpartyID, tc.reference)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestWalletFreeze(t *testing.T) {
	t.Run("freeze active wallet", func(t *testing.T) {
		w := mustWallet(t, "owner-1", "USD", 100)
		require.NoError(t, w.Freeze("compliance-check"))
	})

	t.Run("freeze already frozen wallet", func(t *testing.T) {
		w := mustWallet(t, "owner-1", "USD", 100)
		require.NoError(t, w.Freeze("first-freeze"))
		err := w.Freeze("second-freeze")
		require.ErrorIs(t, err, domain.ErrWalletNotActive)
	})
}

// mustWallet is a test helper that creates a Wallet or fails the test immediately.
func mustWallet(t *testing.T, ownerID, currency string, initialBalance float64) *domain.Wallet {
	t.Helper()
	w, err := domain.NewWallet(ownerID, currency, initialBalance)
	require.NoError(t, err)
	return w
}

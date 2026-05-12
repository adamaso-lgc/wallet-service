package application_test

import (
	"context"
	"testing"

	"github.com/adamaso/wallet-service/internal/application"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newHandler wires up a fresh handler with a fake in-memory repository.
func newHandler(t *testing.T) *application.WalletCommandHandler {
	t.Helper()
	return application.NewWalletCommandHandler(fake.NewWalletRepository())
}

// createWallet is a test helper that creates a wallet and returns its ID.
func createWallet(t *testing.T, h *application.WalletCommandHandler) string {
	t.Helper()
	result, err := h.CreateWallet(context.Background(), application.CreateWalletCommand{
		OwnerID:        "owner-1",
		Currency:       "USD",
		InitialBalance: 100,
	})
	require.NoError(t, err)
	return result.WalletID
}

// --- CreateWallet ---

func TestCreateWallet_Success(t *testing.T) {
	h := newHandler(t)

	result, err := h.CreateWallet(context.Background(), application.CreateWalletCommand{
		OwnerID:        "owner-1",
		Currency:       "USD",
		InitialBalance: 500,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, result.WalletID)
}

func TestCreateWallet_InvalidInput(t *testing.T) {
	h := newHandler(t)

	tests := []struct {
		name    string
		cmd     application.CreateWalletCommand
		wantErr error
	}{
		{
			name:    "missing owner",
			cmd:     application.CreateWalletCommand{Currency: "USD", InitialBalance: 100},
			wantErr: domain.ErrOwnerIDRequired,
		},
		{
			name:    "missing currency",
			cmd:     application.CreateWalletCommand{OwnerID: "owner-1", InitialBalance: 100},
			wantErr: domain.ErrCurrencyRequired,
		},
		{
			name:    "negative balance",
			cmd:     application.CreateWalletCommand{OwnerID: "owner-1", Currency: "USD", InitialBalance: -10},
			wantErr: domain.ErrInvalidAmount,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := h.CreateWallet(context.Background(), tc.cmd)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

// --- Deposit ---

func TestDeposit_Success(t *testing.T) {
	h := newHandler(t)
	id := createWallet(t, h)

	err := h.Deposit(context.Background(), application.DepositCommand{
		WalletID:  id,
		Amount:    50,
		Reference: "top-up",
	})

	require.NoError(t, err)
}

func TestDeposit_WalletNotFound(t *testing.T) {
	h := newHandler(t)

	err := h.Deposit(context.Background(), application.DepositCommand{
		WalletID: "nonexistent",
		Amount:   50,
	})

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestDeposit_OnFrozenWallet(t *testing.T) {
	h := newHandler(t)
	id := createWallet(t, h)
	require.NoError(t, h.FreezeWallet(context.Background(), application.FreezeWalletCommand{WalletID: id}))

	err := h.Deposit(context.Background(), application.DepositCommand{WalletID: id, Amount: 50})

	require.ErrorIs(t, err, domain.ErrWalletNotActive)
}

// --- Withdraw ---

func TestWithdraw_Success(t *testing.T) {
	h := newHandler(t)
	id := createWallet(t, h)

	err := h.Withdraw(context.Background(), application.WithdrawCommand{
		WalletID:  id,
		Amount:    40,
		Reference: "purchase",
	})

	require.NoError(t, err)
}

func TestWithdraw_InsufficientFunds(t *testing.T) {
	h := newHandler(t)
	id := createWallet(t, h) // balance: 100

	err := h.Withdraw(context.Background(), application.WithdrawCommand{
		WalletID: id,
		Amount:   200,
	})

	require.ErrorIs(t, err, domain.ErrInsufficientFunds)
}

func TestWithdraw_WalletNotFound(t *testing.T) {
	h := newHandler(t)

	err := h.Withdraw(context.Background(), application.WithdrawCommand{
		WalletID: "nonexistent",
		Amount:   10,
	})

	require.ErrorIs(t, err, domain.ErrNotFound)
}

// --- Transfer ---

func TestTransfer_Success(t *testing.T) {
	h := newHandler(t)
	sourceID := createWallet(t, h)      // balance: 100
	destinationID := createWallet(t, h) // balance: 100

	err := h.Transfer(context.Background(), application.TransferCommand{
		SourceWalletID:      sourceID,
		DestinationWalletID: destinationID,
		Amount:              60,
		Reference:           "payment",
	})

	require.NoError(t, err)
}

func TestTransfer_SourceNotFound(t *testing.T) {
	h := newHandler(t)
	destinationID := createWallet(t, h)

	err := h.Transfer(context.Background(), application.TransferCommand{
		SourceWalletID:      "nonexistent",
		DestinationWalletID: destinationID,
		Amount:              50,
	})

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestTransfer_DestinationNotFound(t *testing.T) {
	h := newHandler(t)
	sourceID := createWallet(t, h)

	err := h.Transfer(context.Background(), application.TransferCommand{
		SourceWalletID:      sourceID,
		DestinationWalletID: "nonexistent",
		Amount:              50,
	})

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestTransfer_InsufficientFunds(t *testing.T) {
	h := newHandler(t)
	sourceID := createWallet(t, h) // balance: 100
	destinationID := createWallet(t, h)

	err := h.Transfer(context.Background(), application.TransferCommand{
		SourceWalletID:      sourceID,
		DestinationWalletID: destinationID,
		Amount:              200,
	})

	require.ErrorIs(t, err, domain.ErrInsufficientFunds)
}

// --- FreezeWallet ---

func TestFreezeWallet_Success(t *testing.T) {
	h := newHandler(t)
	id := createWallet(t, h)

	err := h.FreezeWallet(context.Background(), application.FreezeWalletCommand{
		WalletID:  id,
		Reference: "compliance",
	})

	require.NoError(t, err)
}

func TestFreezeWallet_AlreadyFrozen(t *testing.T) {
	h := newHandler(t)
	id := createWallet(t, h)
	require.NoError(t, h.FreezeWallet(context.Background(), application.FreezeWalletCommand{WalletID: id}))

	err := h.FreezeWallet(context.Background(), application.FreezeWalletCommand{WalletID: id})

	require.ErrorIs(t, err, domain.ErrWalletNotActive)
}

func TestFreezeWallet_NotFound(t *testing.T) {
	h := newHandler(t)

	err := h.FreezeWallet(context.Background(), application.FreezeWalletCommand{WalletID: "nonexistent"})

	require.ErrorIs(t, err, domain.ErrNotFound)
}

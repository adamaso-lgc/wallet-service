package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adamaso/wallet-service/internal/application/command"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/fake"
)

func TestCreateWallet_Success(t *testing.T) {
	h := command.NewCreateWalletHandler(fake.NewWalletRepository())

	result, err := h.Handle(context.Background(), command.CreateWallet{
		OwnerID:        "owner-1",
		Currency:       "USD",
		InitialBalance: 500,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, result.WalletID)
}

func TestCreateWallet_InvalidInput(t *testing.T) {
	h := command.NewCreateWalletHandler(fake.NewWalletRepository())

	tests := []struct {
		name    string
		cmd     command.CreateWallet
		wantErr error
	}{
		{
			name:    "missing owner",
			cmd:     command.CreateWallet{Currency: "USD", InitialBalance: 100},
			wantErr: domain.ErrOwnerIDRequired,
		},
		{
			name:    "missing currency",
			cmd:     command.CreateWallet{OwnerID: "owner-1", InitialBalance: 100},
			wantErr: domain.ErrCurrencyRequired,
		},
		{
			name:    "negative balance",
			cmd:     command.CreateWallet{OwnerID: "owner-1", Currency: "USD", InitialBalance: -10},
			wantErr: domain.ErrInvalidAmount,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := h.Handle(context.Background(), tc.cmd)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

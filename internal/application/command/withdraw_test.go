package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/adamaso/wallet-service/internal/application/command"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/fake"
)

func TestWithdraw_Success(t *testing.T) {
	repo := fake.NewWalletRepository()
	id := mustCreateWallet(t, repo)

	err := command.NewWithdrawHandler(repo).Handle(context.Background(), command.Withdraw{
		WalletID:  id,
		Amount:    40,
		Reference: "purchase",
	})

	require.NoError(t, err)
}

func TestWithdraw_InsufficientFunds(t *testing.T) {
	repo := fake.NewWalletRepository()
	id := mustCreateWallet(t, repo) // balance: 100

	err := command.NewWithdrawHandler(repo).Handle(
		context.Background(),
		command.Withdraw{WalletID: id, Amount: 200},
	)

	require.ErrorIs(t, err, domain.ErrInsufficientFunds)
}

func TestWithdraw_WalletNotFound(t *testing.T) {
	err := command.NewWithdrawHandler(fake.NewWalletRepository()).Handle(
		context.Background(),
		command.Withdraw{WalletID: "nonexistent", Amount: 10},
	)

	require.ErrorIs(t, err, domain.ErrNotFound)
}

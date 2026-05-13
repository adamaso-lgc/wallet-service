package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/adamaso/wallet-service/internal/application/command"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/fake"
)

func TestDeposit_Success(t *testing.T) {
	repo := fake.NewWalletRepository()
	id := mustCreateWallet(t, repo)

	err := command.NewDepositHandler(repo).Handle(context.Background(), command.Deposit{
		WalletID:  id,
		Amount:    50,
		Reference: "top-up",
	})

	require.NoError(t, err)
}

func TestDeposit_WalletNotFound(t *testing.T) {
	err := command.NewDepositHandler(fake.NewWalletRepository()).Handle(
		context.Background(),
		command.Deposit{WalletID: "nonexistent", Amount: 50},
	)

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestDeposit_OnFrozenWallet(t *testing.T) {
	repo := fake.NewWalletRepository()
	id := mustCreateWallet(t, repo)
	require.NoError(t, command.NewFreezeWalletHandler(repo).Handle(
		context.Background(), command.FreezeWallet{WalletID: id},
	))

	err := command.NewDepositHandler(repo).Handle(
		context.Background(),
		command.Deposit{WalletID: id, Amount: 50},
	)

	require.ErrorIs(t, err, domain.ErrWalletNotActive)
}

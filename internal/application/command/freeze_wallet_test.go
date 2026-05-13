package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/adamaso/wallet-service/internal/application/command"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/fake"
)

func TestFreezeWallet_Success(t *testing.T) {
	repo := fake.NewWalletRepository()
	id := mustCreateWallet(t, repo)

	err := command.NewFreezeWalletHandler(repo).Handle(context.Background(), command.FreezeWallet{
		WalletID:  id,
		Reference: "compliance",
	})

	require.NoError(t, err)
}

func TestFreezeWallet_AlreadyFrozen(t *testing.T) {
	repo := fake.NewWalletRepository()
	id := mustCreateWallet(t, repo)
	require.NoError(t, command.NewFreezeWalletHandler(repo).Handle(
		context.Background(), command.FreezeWallet{WalletID: id},
	))

	err := command.NewFreezeWalletHandler(repo).Handle(
		context.Background(),
		command.FreezeWallet{WalletID: id},
	)

	require.ErrorIs(t, err, domain.ErrWalletNotActive)
}

func TestFreezeWallet_NotFound(t *testing.T) {
	err := command.NewFreezeWalletHandler(fake.NewWalletRepository()).Handle(
		context.Background(),
		command.FreezeWallet{WalletID: "nonexistent"},
	)

	require.ErrorIs(t, err, domain.ErrNotFound)
}

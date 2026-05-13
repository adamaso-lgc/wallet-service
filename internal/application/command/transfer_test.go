package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/adamaso/wallet-service/internal/application/command"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/fake"
)

func TestTransfer_Success(t *testing.T) {
	repo := fake.NewWalletRepository()
	sourceID := mustCreateWallet(t, repo)
	destinationID := mustCreateWallet(t, repo)

	err := command.NewTransferHandler(repo).Handle(context.Background(), command.Transfer{
		SourceWalletID:      sourceID,
		DestinationWalletID: destinationID,
		Amount:              60,
		Reference:           "payment",
	})

	require.NoError(t, err)
}

func TestTransfer_SourceNotFound(t *testing.T) {
	repo := fake.NewWalletRepository()
	destinationID := mustCreateWallet(t, repo)

	err := command.NewTransferHandler(repo).Handle(context.Background(), command.Transfer{
		SourceWalletID:      "nonexistent",
		DestinationWalletID: destinationID,
		Amount:              50,
	})

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestTransfer_DestinationNotFound(t *testing.T) {
	repo := fake.NewWalletRepository()
	sourceID := mustCreateWallet(t, repo)

	err := command.NewTransferHandler(repo).Handle(context.Background(), command.Transfer{
		SourceWalletID:      sourceID,
		DestinationWalletID: "nonexistent",
		Amount:              50,
	})

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestTransfer_InsufficientFunds(t *testing.T) {
	repo := fake.NewWalletRepository()
	sourceID := mustCreateWallet(t, repo) // balance: 100
	destinationID := mustCreateWallet(t, repo)

	err := command.NewTransferHandler(repo).Handle(context.Background(), command.Transfer{
		SourceWalletID:      sourceID,
		DestinationWalletID: destinationID,
		Amount:              200,
	})

	require.ErrorIs(t, err, domain.ErrInsufficientFunds)
}

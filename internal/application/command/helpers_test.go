package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/adamaso/wallet-service/internal/application/command"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/fake"
)

// mustCreateWallet creates a wallet with a 100 USD balance and returns its ID.
// It is shared across command test files in this package.
func mustCreateWallet(t *testing.T, repo domain.WalletRepository) string {
	t.Helper()
	result, err := command.NewCreateWalletHandler(repo).Handle(
		context.Background(),
		command.CreateWallet{OwnerID: "owner-1", Currency: "USD", InitialBalance: 100},
	)
	require.NoError(t, err)
	return result.WalletID
}

// newRepo returns a fresh fake repository for each test.
func newRepo() *fake.WalletRepository {
	return fake.NewWalletRepository()
}

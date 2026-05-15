package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/eventstore"
	"github.com/adamaso/wallet-service/internal/infrastructure/projection"
	"github.com/adamaso/wallet-service/internal/service"
)

// seedWallet creates a wallet through the domain and saves it via the repository,
// without going through the service under test.
func seedWallet(t *testing.T, ownerID, currency string, balance float64) *domain.Wallet {
	t.Helper()
	wallet, err := domain.NewWallet(ownerID, currency, balance)
	require.NoError(t, err)
	require.NoError(t, newRepo(t).Save(context.Background(), wallet))
	return wallet
}

// seedFrozenWallet creates and immediately freezes a wallet.
func seedFrozenWallet(t *testing.T, ownerID, currency string, balance float64) *domain.Wallet {
	t.Helper()
	wallet := seedWallet(t, ownerID, currency, balance)
	require.NoError(t, wallet.Freeze("test-freeze"))
	require.NoError(t, newRepo(t).Save(context.Background(), wallet))
	return wallet
}

// getView fetches the wallet projection via the service — used to assert state after commands.
func getView(t *testing.T, svc *service.Service, walletID string) *walletv1.WalletResponse {
	t.Helper()
	resp, err := svc.GetWallet(context.Background(), &walletv1.GetWalletRequest{WalletId: walletID})
	require.NoError(t, err)
	return resp
}

func newRepo(t *testing.T) *eventstore.WalletRepository {
	t.Helper()
	return eventstore.NewWalletRepository(testPool, projection.NewWalletProjector())
}

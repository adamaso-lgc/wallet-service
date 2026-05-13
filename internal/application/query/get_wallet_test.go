package query_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adamaso/wallet-service/internal/application/query"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/fake"
	"github.com/adamaso/wallet-service/internal/projection"
)

func walletView(id, ownerID string) *projection.WalletView {
	return &projection.WalletView{
		ID: id, OwnerID: ownerID,
		Balance: 100.0, Currency: "USD", Status: "active",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

func TestGetWallet_ReturnsView(t *testing.T) {
	view := walletView("wallet-1", "owner-1")
	h := query.NewGetWalletHandler(fake.NewWalletViewStore(view))

	result, err := h.Handle(context.Background(), query.GetWalletQuery{ID: "wallet-1"})

	require.NoError(t, err)
	assert.Equal(t, view.ID, result.ID)
	assert.Equal(t, view.OwnerID, result.OwnerID)
	assert.Equal(t, view.Balance, result.Balance)
	assert.Equal(t, view.Currency, result.Currency)
}

func TestGetWallet_NotFound(t *testing.T) {
	h := query.NewGetWalletHandler(fake.NewWalletViewStore())

	_, err := h.Handle(context.Background(), query.GetWalletQuery{ID: "nonexistent"})

	require.ErrorIs(t, err, domain.ErrNotFound)
}

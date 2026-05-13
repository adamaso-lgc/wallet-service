package query_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adamaso/wallet-service/internal/application/query"
	"github.com/adamaso/wallet-service/internal/infrastructure/fake"
)

func TestListWalletsByOwner_ReturnsAllForOwner(t *testing.T) {
	store := fake.NewWalletViewStore(
		walletView("wallet-1", "owner-1"),
		walletView("wallet-2", "owner-1"),
		walletView("wallet-3", "owner-2"),
	)
	h := query.NewListWalletsByOwnerHandler(store)

	results, err := h.Handle(context.Background(), query.ListWalletsByOwnerQuery{OwnerID: "owner-1"})

	require.NoError(t, err)
	assert.Len(t, results, 2)
	for _, r := range results {
		assert.Equal(t, "owner-1", r.OwnerID)
	}
}

func TestListWalletsByOwner_UnknownOwnerReturnsEmpty(t *testing.T) {
	store := fake.NewWalletViewStore(walletView("wallet-1", "owner-1"))
	h := query.NewListWalletsByOwnerHandler(store)

	results, err := h.Handle(context.Background(), query.ListWalletsByOwnerQuery{OwnerID: "unknown"})

	require.NoError(t, err)
	assert.Empty(t, results)
}

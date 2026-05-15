package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
)

func TestListWalletsByOwner_ReturnsOwnersWallets(t *testing.T) {
	svc := newTestService(t)
	w1 := seedWallet(t, "owner-1", "EUR", 100)
	w2 := seedWallet(t, "owner-1", "USD", 200)
	seedWallet(t, "owner-2", "EUR", 50)

	resp, err := svc.ListWalletsByOwner(context.Background(), &walletv1.ListWalletsByOwnerRequest{
		OwnerId: "owner-1",
	})

	require.NoError(t, err)
	require.Len(t, resp.Wallets, 2)

	ids := []string{resp.Wallets[0].WalletId, resp.Wallets[1].WalletId}
	assert.ElementsMatch(t, []string{w1.GetID(), w2.GetID()}, ids)
	for _, w := range resp.Wallets {
		assert.Equal(t, "owner-1", w.OwnerId)
	}
}

func TestListWalletsByOwner_EmptyForUnknownOwner(t *testing.T) {
	svc := newTestService(t)
	seedWallet(t, "owner-1", "EUR", 100)

	resp, err := svc.ListWalletsByOwner(context.Background(), &walletv1.ListWalletsByOwnerRequest{
		OwnerId: "unknown-owner",
	})

	require.NoError(t, err)
	assert.Empty(t, resp.Wallets)
}

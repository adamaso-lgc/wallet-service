package query

import (
	"context"

	"github.com/adamaso/wallet-service/internal/projection"
)

type ListWalletsByOwnerQuery struct {
	OwnerID string
}

type ListWalletsByOwnerHandler struct {
	store projection.WalletStore
}

func NewListWalletsByOwnerHandler(store projection.WalletStore) *ListWalletsByOwnerHandler {
	return &ListWalletsByOwnerHandler{store: store}
}

func (h *ListWalletsByOwnerHandler) Handle(ctx context.Context, q ListWalletsByOwnerQuery) ([]*projection.WalletView, error) {
	return h.store.ListWalletsByOwner(ctx, q.OwnerID)
}

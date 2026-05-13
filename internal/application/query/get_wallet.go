package query

import (
	"context"

	"github.com/adamaso/wallet-service/internal/projection"
)

type GetWalletQuery struct {
	ID string
}

type GetWalletHandler struct {
	store projection.WalletStore
}

func NewGetWalletHandler(store projection.WalletStore) *GetWalletHandler {
	return &GetWalletHandler{store: store}
}

func (h *GetWalletHandler) Handle(ctx context.Context, q GetWalletQuery) (*projection.WalletView, error) {
	return h.store.GetWallet(ctx, q.ID)
}

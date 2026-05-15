package query

import (
	"context"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application/common"
	"github.com/adamaso/wallet-service/internal/infrastructure/projection"
)

type ListWalletsByOwnerHandler struct {
	store projection.Repository
}

func NewListWalletsByOwnerHandler(store projection.Repository) *ListWalletsByOwnerHandler {
	return &ListWalletsByOwnerHandler{store: store}
}

func (h *ListWalletsByOwnerHandler) Handle(ctx context.Context, req *walletv1.ListWalletsByOwnerRequest) (*walletv1.ListWalletsByOwnerResponse, error) {
	views, err := h.store.ListWalletsByOwner(ctx, req.OwnerId)
	if err != nil {
		return nil, common.ToGRPCError(err)
	}
	wallets := make([]*walletv1.WalletResponse, len(views))
	for i, v := range views {
		wallets[i] = common.ToWalletResponse(v)
	}
	return &walletv1.ListWalletsByOwnerResponse{Wallets: wallets}, nil
}

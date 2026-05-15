package query

import (
	"context"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application/common"
	"github.com/adamaso/wallet-service/internal/infrastructure/projection"
)

type GetWalletHandler struct {
	store projection.Repository
}

func NewGetWalletHandler(store projection.Repository) *GetWalletHandler {
	return &GetWalletHandler{store: store}
}

func (h *GetWalletHandler) Handle(ctx context.Context, req *walletv1.GetWalletRequest) (*walletv1.WalletResponse, error) {
	v, err := h.store.GetWallet(ctx, req.WalletId)
	if err != nil {
		return nil, common.ToGRPCError(err)
	}
	return common.ToWalletResponse(v), nil
}

package command

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application/common"
	"github.com/adamaso/wallet-service/internal/domain"
)

type FreezeWalletHandler struct {
	repo domain.WalletRepository
}

func NewFreezeWalletHandler(repo domain.WalletRepository) *FreezeWalletHandler {
	return &FreezeWalletHandler{repo: repo}
}

func (h *FreezeWalletHandler) Handle(ctx context.Context, req *walletv1.FreezeWalletRequest) (*emptypb.Empty, error) {
	wallet, err := h.repo.Get(ctx, req.WalletId)
	if err != nil {
		return nil, common.ToGRPCError(err)
	}
	if err := wallet.Freeze(req.Reference); err != nil {
		return nil, common.ToGRPCError(err)
	}
	if err := h.repo.Save(ctx, wallet); err != nil {
		return nil, common.ToGRPCError(err)
	}
	return &emptypb.Empty{}, nil
}

package command

import (
	"context"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application/common"
	"github.com/adamaso/wallet-service/internal/domain"
)

type CreateWalletHandler struct {
	repo domain.WalletRepository
}

func NewCreateWalletHandler(repo domain.WalletRepository) *CreateWalletHandler {
	return &CreateWalletHandler{repo: repo}
}

func (h *CreateWalletHandler) Handle(ctx context.Context, req *walletv1.CreateWalletRequest) (*walletv1.CreateWalletResponse, error) {
	wallet, err := domain.NewWallet(req.OwnerId, req.Currency, req.InitialBalance)
	if err != nil {
		return nil, common.ToGRPCError(err)
	}
	if err := h.repo.Save(ctx, wallet); err != nil {
		return nil, common.ToGRPCError(err)
	}

	return &walletv1.CreateWalletResponse{WalletId: wallet.GetID()}, nil
}

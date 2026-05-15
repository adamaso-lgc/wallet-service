package command

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application/common"
	"github.com/adamaso/wallet-service/internal/domain"
)

type WithdrawHandler struct {
	repo domain.WalletRepository
}

func NewWithdrawHandler(repo domain.WalletRepository) *WithdrawHandler {
	return &WithdrawHandler{repo: repo}
}

func (h *WithdrawHandler) Handle(ctx context.Context, req *walletv1.WithdrawRequest) (*emptypb.Empty, error) {
	wallet, err := h.repo.Get(ctx, req.WalletId)
	if err != nil {
		return nil, common.ToGRPCError(err)
	}
	if err := wallet.Withdraw(req.Amount, req.Reference); err != nil {
		return nil, common.ToGRPCError(err)
	}
	if err := h.repo.Save(ctx, wallet); err != nil {
		return nil, common.ToGRPCError(err)
	}
	return &emptypb.Empty{}, nil
}

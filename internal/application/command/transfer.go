package command

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application/common"
	"github.com/adamaso/wallet-service/internal/domain"
)

// TransferHandler moves funds between two wallets atomically via SaveAll.
// Both wallets are saved in a single transaction — either both succeed or
// neither does, preventing partial state.
type TransferHandler struct {
	repo domain.WalletRepository
}

func NewTransferHandler(repo domain.WalletRepository) *TransferHandler {
	return &TransferHandler{repo: repo}
}

func (h *TransferHandler) Handle(ctx context.Context, req *walletv1.TransferRequest) (*emptypb.Empty, error) {
	if req.SourceWalletId == req.DestinationWalletId {
		return nil, status.Error(codes.InvalidArgument, "cannot transfer to self")
	}

	source, err := h.repo.Get(ctx, req.SourceWalletId)
	if err != nil {
		return nil, common.ToGRPCError(err)
	}
	destination, err := h.repo.Get(ctx, req.DestinationWalletId)
	if err != nil {
		return nil, common.ToGRPCError(err)
	}
	if err := source.DebitForTransfer(req.Amount, req.DestinationWalletId, req.Reference); err != nil {
		return nil, common.ToGRPCError(err)
	}
	if err := destination.CreditForTransfer(req.Amount, req.SourceWalletId, req.Reference); err != nil {
		return nil, common.ToGRPCError(err)
	}
	if err := h.repo.SaveAll(ctx, source, destination); err != nil {
		return nil, common.ToGRPCError(err)
	}
	return &emptypb.Empty{}, nil
}

package grpcserver

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application"
	"github.com/adamaso/wallet-service/internal/application/command"
	"github.com/adamaso/wallet-service/internal/application/query"
)

type WalletHandler struct {
	walletv1.UnimplementedWalletServiceServer
	app *application.Application
}

var _ walletv1.WalletServiceServer = (*WalletHandler)(nil)

func NewWalletHandler(app *application.Application) *WalletHandler {
	return &WalletHandler{app: app}
}

func (h *WalletHandler) CreateWallet(ctx context.Context, req *walletv1.CreateWalletRequest) (*walletv1.CreateWalletResponse, error) {
	result, err := h.app.Commands.CreateWallet.Handle(ctx, command.CreateWallet{
		OwnerID:        req.OwnerId,
		Currency:       req.Currency,
		InitialBalance: req.InitialBalance,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &walletv1.CreateWalletResponse{WalletId: result.WalletID}, nil
}

func (h *WalletHandler) Deposit(ctx context.Context, req *walletv1.DepositRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.Deposit.Handle(ctx, command.Deposit{
		WalletID:  req.WalletId,
		Amount:    req.Amount,
		Reference: req.Reference,
	})
	return &emptypb.Empty{}, toGRPCError(err)
}

func (h *WalletHandler) Withdraw(ctx context.Context, req *walletv1.WithdrawRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.Withdraw.Handle(ctx, command.Withdraw{
		WalletID:  req.WalletId,
		Amount:    req.Amount,
		Reference: req.Reference,
	})
	return &emptypb.Empty{}, toGRPCError(err)
}

func (h *WalletHandler) Transfer(ctx context.Context, req *walletv1.TransferRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.Transfer.Handle(ctx, command.Transfer{
		SourceWalletID:      req.SourceWalletId,
		DestinationWalletID: req.DestinationWalletId,
		Amount:              req.Amount,
		Reference:           req.Reference,
	})
	return &emptypb.Empty{}, toGRPCError(err)
}

func (h *WalletHandler) FreezeWallet(ctx context.Context, req *walletv1.FreezeWalletRequest) (*emptypb.Empty, error) {
	err := h.app.Commands.FreezeWallet.Handle(ctx, command.FreezeWallet{
		WalletID:  req.WalletId,
		Reference: req.Reference,
	})
	return &emptypb.Empty{}, toGRPCError(err)
}

func (h *WalletHandler) GetWallet(ctx context.Context, req *walletv1.GetWalletRequest) (*walletv1.WalletResponse, error) {
	view, err := h.app.Queries.GetWallet.Handle(ctx, query.GetWalletQuery{ID: req.WalletId})
	if err != nil {
		return nil, toGRPCError(err)
	}
	return toWalletResponse(view), nil
}

func (h *WalletHandler) ListWalletsByOwner(ctx context.Context, req *walletv1.ListWalletsByOwnerRequest) (*walletv1.ListWalletsByOwnerResponse, error) {
	views, err := h.app.Queries.ListWalletsByOwner.Handle(ctx, query.ListWalletsByOwnerQuery{OwnerID: req.OwnerId})
	if err != nil {
		return nil, toGRPCError(err)
	}
	resp := &walletv1.ListWalletsByOwnerResponse{
		Wallets: make([]*walletv1.WalletResponse, len(views)),
	}
	for i, v := range views {
		resp.Wallets[i] = toWalletResponse(v)
	}
	return resp, nil
}

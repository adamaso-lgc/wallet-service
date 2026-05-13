package grpcserver

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application/command"
)

func (s *Server) CreateWallet(ctx context.Context, req *walletv1.CreateWalletRequest) (*walletv1.CreateWalletResponse, error) {
	result, err := s.app.Commands.CreateWallet.Handle(ctx, command.CreateWallet{
		OwnerID:        req.OwnerId,
		Currency:       req.Currency,
		InitialBalance: req.InitialBalance,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &walletv1.CreateWalletResponse{WalletId: result.WalletID}, nil
}

func (s *Server) Deposit(ctx context.Context, req *walletv1.DepositRequest) (*emptypb.Empty, error) {
	err := s.app.Commands.Deposit.Handle(ctx, command.Deposit{
		WalletID:  req.WalletId,
		Amount:    req.Amount,
		Reference: req.Reference,
	})
	return &emptypb.Empty{}, toGRPCError(err)
}

func (s *Server) Withdraw(ctx context.Context, req *walletv1.WithdrawRequest) (*emptypb.Empty, error) {
	err := s.app.Commands.Withdraw.Handle(ctx, command.Withdraw{
		WalletID:  req.WalletId,
		Amount:    req.Amount,
		Reference: req.Reference,
	})
	return &emptypb.Empty{}, toGRPCError(err)
}

func (s *Server) Transfer(ctx context.Context, req *walletv1.TransferRequest) (*emptypb.Empty, error) {
	err := s.app.Commands.Transfer.Handle(ctx, command.Transfer{
		SourceWalletID:      req.SourceWalletId,
		DestinationWalletID: req.DestinationWalletId,
		Amount:              req.Amount,
		Reference:           req.Reference,
	})
	return &emptypb.Empty{}, toGRPCError(err)
}

func (s *Server) FreezeWallet(ctx context.Context, req *walletv1.FreezeWalletRequest) (*emptypb.Empty, error) {
	err := s.app.Commands.FreezeWallet.Handle(ctx, command.FreezeWallet{
		WalletID:  req.WalletId,
		Reference: req.Reference,
	})
	return &emptypb.Empty{}, toGRPCError(err)
}

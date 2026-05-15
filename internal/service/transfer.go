package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
)

// Transfer moves funds between two wallets atomically via SaveAll.
// Both wallets are saved in a single transaction — either both succeed or
// neither does, preventing partial state.
func (s *Service) Transfer(ctx context.Context, req *walletv1.TransferRequest) (*emptypb.Empty, error) {
	if req.SourceWalletId == req.DestinationWalletId {
		return nil, status.Error(codes.InvalidArgument, "cannot transfer to self")
	}

	source, err := s.repo.Get(ctx, req.SourceWalletId)
	if err != nil {
		return nil, toGRPCError(err)
	}
	destination, err := s.repo.Get(ctx, req.DestinationWalletId)
	if err != nil {
		return nil, toGRPCError(err)
	}
	if err := source.DebitForTransfer(req.Amount, req.DestinationWalletId, req.Reference); err != nil {
		return nil, toGRPCError(err)
	}
	if err := destination.CreditForTransfer(req.Amount, req.SourceWalletId, req.Reference); err != nil {
		return nil, toGRPCError(err)
	}
	if err := s.repo.SaveAll(ctx, source, destination); err != nil {
		return nil, toGRPCError(err)
	}
	return &emptypb.Empty{}, nil
}

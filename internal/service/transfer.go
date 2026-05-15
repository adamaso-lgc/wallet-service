package service

import (
	"context"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/domain"
)

// Transfer moves funds between two wallets atomically via SaveAll.
func (s *Service) Transfer(ctx context.Context, req *walletv1.TransferRequest) (*emptypb.Empty, error) {
	if req.SourceWalletId == req.DestinationWalletId {
		return nil, status.Error(codes.InvalidArgument, "cannot transfer to self")
	}

	var source, destination *domain.Wallet

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		source, err = s.repo.Get(ctx, req.SourceWalletId)
		return err
	})
	g.Go(func() error {
		var err error
		destination, err = s.repo.Get(ctx, req.DestinationWalletId)
		return err
	})
	if err := g.Wait(); err != nil {
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

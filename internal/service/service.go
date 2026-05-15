package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/projection"
)

type Service struct {
	walletv1.UnimplementedWalletServiceServer
	repo  domain.WalletRepository
	store projection.Repository
}

var _ walletv1.WalletServiceServer = (*Service)(nil)

func NewService(repo domain.WalletRepository, store projection.Repository) *Service {
	return &Service{repo: repo, store: store}
}

// applyAndSave fetches a wallet, runs fn against it, and persists the result.
func (s *Service) applyAndSave(ctx context.Context, walletID string, fn func(*domain.Wallet) error) (*emptypb.Empty, error) {
	wallet, err := s.repo.Get(ctx, walletID)
	if err != nil {
		return nil, toGRPCError(err)
	}
	if err := fn(wallet); err != nil {
		return nil, toGRPCError(err)
	}
	if err := s.repo.Save(ctx, wallet); err != nil {
		return nil, toGRPCError(err)
	}
	return &emptypb.Empty{}, nil
}

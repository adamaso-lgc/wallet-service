package service

import (
	"context"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/domain"
)

func (s *Service) CreateWallet(ctx context.Context, req *walletv1.CreateWalletRequest) (*walletv1.CreateWalletResponse, error) {
	wallet, err := domain.NewWallet(req.OwnerId, req.Currency, req.InitialBalance)
	if err != nil {
		return nil, toGRPCError(err)
	}
	if err := s.repo.Save(ctx, wallet); err != nil {
		return nil, toGRPCError(err)
	}
	return &walletv1.CreateWalletResponse{WalletId: wallet.GetID()}, nil
}

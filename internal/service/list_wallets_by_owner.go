package service

import (
	"context"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
)

func (s *Service) ListWalletsByOwner(ctx context.Context, req *walletv1.ListWalletsByOwnerRequest) (*walletv1.ListWalletsByOwnerResponse, error) {
	views, err := s.store.ListWalletsByOwner(ctx, req.OwnerId)
	if err != nil {
		return nil, toGRPCError(err)
	}
	wallets := make([]*walletv1.WalletResponse, len(views))
	for i, v := range views {
		wallets[i] = toWalletResponse(v)
	}
	return &walletv1.ListWalletsByOwnerResponse{Wallets: wallets}, nil
}

package service

import (
	"context"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
)

func (s *Service) GetWallet(ctx context.Context, req *walletv1.GetWalletRequest) (*walletv1.WalletResponse, error) {
	v, err := s.store.GetWallet(ctx, req.WalletId)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return toWalletResponse(v), nil
}

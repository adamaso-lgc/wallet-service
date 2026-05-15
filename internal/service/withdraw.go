package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/domain"
)

func (s *Service) Withdraw(ctx context.Context, req *walletv1.WithdrawRequest) (*emptypb.Empty, error) {
	return s.applyAndSave(ctx, req.WalletId, func(w *domain.Wallet) error {
		return w.Withdraw(req.Amount, req.Reference)
	})
}

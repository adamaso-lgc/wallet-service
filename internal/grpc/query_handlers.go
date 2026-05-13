package grpcserver

import (
	"context"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application/query"
)

func (s *Server) GetWallet(ctx context.Context, req *walletv1.GetWalletRequest) (*walletv1.WalletResponse, error) {
	view, err := s.app.Queries.GetWallet.Handle(ctx, query.GetWalletQuery{ID: req.WalletId})
	if err != nil {
		return nil, toGRPCError(err)
	}
	return toWalletResponse(view), nil
}

func (s *Server) ListWalletsByOwner(ctx context.Context, req *walletv1.ListWalletsByOwnerRequest) (*walletv1.ListWalletsByOwnerResponse, error) {
	views, err := s.app.Queries.ListWalletsByOwner.Handle(ctx, query.ListWalletsByOwnerQuery{OwnerID: req.OwnerId})
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

package service

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/infrastructure/projection"
)

func toWalletResponse(v *projection.WalletView) *walletv1.WalletResponse {
	return &walletv1.WalletResponse{
		WalletId:  v.ID,
		OwnerId:   v.OwnerID,
		Balance:   v.Balance,
		Currency:  v.Currency,
		Status:    v.Status,
		CreatedAt: timestamppb.New(v.CreatedAt),
		UpdatedAt: timestamppb.New(v.UpdatedAt),
	}
}

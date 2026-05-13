package grpcserver

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/projection"
)

// toWalletResponse converts a projection.WalletView to the proto response type.
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

// toGRPCError maps domain sentinel errors to the correct gRPC status code.
// nil is returned as nil so callers can write: return x, toGRPCError(err)
func toGRPCError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrOwnerIDRequired),
		errors.Is(err, domain.ErrCurrencyRequired),
		errors.Is(err, domain.ErrInvalidAmount),
		errors.Is(err, domain.ErrCurrencyMismatch):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrInsufficientFunds),
		errors.Is(err, domain.ErrWalletNotActive):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrConcurrentModification):
		return status.Error(codes.Aborted, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}

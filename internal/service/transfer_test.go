package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
)

func TestTransfer_Success(t *testing.T) {
	svc := newTestService(t)
	source := seedWallet(t, "owner-1", "EUR", 200)
	destination := seedWallet(t, "owner-2", "EUR", 10)

	_, err := svc.Transfer(context.Background(), &walletv1.TransferRequest{
		SourceWalletId:      source.GetID(),
		DestinationWalletId: destination.GetID(),
		Amount:              75,
		Reference:           "payment",
	})

	require.NoError(t, err)
	assert.Equal(t, float64(125), getView(t, svc, source.GetID()).Balance)
	assert.Equal(t, float64(85), getView(t, svc, destination.GetID()).Balance)
}

func TestTransfer_ToSelf(t *testing.T) {
	svc := newTestService(t)
	wallet := seedWallet(t, "owner-1", "EUR", 100)

	_, err := svc.Transfer(context.Background(), &walletv1.TransferRequest{
		SourceWalletId:      wallet.GetID(),
		DestinationWalletId: wallet.GetID(),
		Amount:              10,
		Reference:           "self",
	})

	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestTransfer_SourceNotFound(t *testing.T) {
	svc := newTestService(t)
	destination := seedWallet(t, "owner-2", "EUR", 10)

	_, err := svc.Transfer(context.Background(), &walletv1.TransferRequest{
		SourceWalletId:      "00000000-0000-0000-0000-000000000000",
		DestinationWalletId: destination.GetID(),
		Amount:              10,
		Reference:           "payment",
	})

	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestTransfer_DestinationNotFound(t *testing.T) {
	svc := newTestService(t)
	source := seedWallet(t, "owner-1", "EUR", 100)

	_, err := svc.Transfer(context.Background(), &walletv1.TransferRequest{
		SourceWalletId:      source.GetID(),
		DestinationWalletId: "00000000-0000-0000-0000-000000000000",
		Amount:              10,
		Reference:           "payment",
	})

	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestTransfer_InsufficientFunds(t *testing.T) {
	svc := newTestService(t)
	source := seedWallet(t, "owner-1", "EUR", 50)
	destination := seedWallet(t, "owner-2", "EUR", 10)

	_, err := svc.Transfer(context.Background(), &walletv1.TransferRequest{
		SourceWalletId:      source.GetID(),
		DestinationWalletId: destination.GetID(),
		Amount:              100,
		Reference:           "payment",
	})

	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestTransfer_FrozenSource(t *testing.T) {
	svc := newTestService(t)
	source := seedFrozenWallet(t, "owner-1", "EUR", 100)
	destination := seedWallet(t, "owner-2", "EUR", 10)

	_, err := svc.Transfer(context.Background(), &walletv1.TransferRequest{
		SourceWalletId:      source.GetID(),
		DestinationWalletId: destination.GetID(),
		Amount:              50,
		Reference:           "payment",
	})

	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestTransfer_FrozenDestination(t *testing.T) {
	svc := newTestService(t)
	source := seedWallet(t, "owner-1", "EUR", 100)
	destination := seedFrozenWallet(t, "owner-2", "EUR", 10)

	_, err := svc.Transfer(context.Background(), &walletv1.TransferRequest{
		SourceWalletId:      source.GetID(),
		DestinationWalletId: destination.GetID(),
		Amount:              50,
		Reference:           "payment",
	})

	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

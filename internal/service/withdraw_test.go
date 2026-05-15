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

func TestWithdraw_Success(t *testing.T) {
	svc := newTestService(t)
	wallet := seedWallet(t, "owner-1", "EUR", 100)

	_, err := svc.Withdraw(context.Background(), &walletv1.WithdrawRequest{
		WalletId:  wallet.GetID(),
		Amount:    40,
		Reference: "purchase",
	})

	require.NoError(t, err)
	assert.Equal(t, float64(60), getView(t, svc, wallet.GetID()).Balance)
}

func TestWithdraw_WalletNotFound(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.Withdraw(context.Background(), &walletv1.WithdrawRequest{
		WalletId:  "00000000-0000-0000-0000-000000000000",
		Amount:    10,
		Reference: "purchase",
	})

	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestWithdraw_InsufficientFunds(t *testing.T) {
	svc := newTestService(t)
	wallet := seedWallet(t, "owner-1", "EUR", 50)

	_, err := svc.Withdraw(context.Background(), &walletv1.WithdrawRequest{
		WalletId:  wallet.GetID(),
		Amount:    100,
		Reference: "purchase",
	})

	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestWithdraw_InvalidAmount(t *testing.T) {
	svc := newTestService(t)
	wallet := seedWallet(t, "owner-1", "EUR", 100)

	_, err := svc.Withdraw(context.Background(), &walletv1.WithdrawRequest{
		WalletId:  wallet.GetID(),
		Amount:    -10,
		Reference: "purchase",
	})

	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestWithdraw_FrozenWallet(t *testing.T) {
	svc := newTestService(t)
	wallet := seedFrozenWallet(t, "owner-1", "EUR", 100)

	_, err := svc.Withdraw(context.Background(), &walletv1.WithdrawRequest{
		WalletId:  wallet.GetID(),
		Amount:    10,
		Reference: "purchase",
	})

	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

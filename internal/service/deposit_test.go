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

func TestDeposit_Success(t *testing.T) {
	svc := newTestService(t)
	wallet := seedWallet(t, "owner-1", "EUR", 100)

	_, err := svc.Deposit(context.Background(), &walletv1.DepositRequest{
		WalletId:  wallet.GetID(),
		Amount:    50,
		Reference: "top-up",
	})

	require.NoError(t, err)
	assert.Equal(t, float64(150), getView(t, svc, wallet.GetID()).Balance)
}

func TestDeposit_WalletNotFound(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.Deposit(context.Background(), &walletv1.DepositRequest{
		WalletId:  "00000000-0000-0000-0000-000000000000",
		Amount:    50,
		Reference: "top-up",
	})

	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestDeposit_InvalidAmount(t *testing.T) {
	svc := newTestService(t)
	wallet := seedWallet(t, "owner-1", "EUR", 100)

	_, err := svc.Deposit(context.Background(), &walletv1.DepositRequest{
		WalletId:  wallet.GetID(),
		Amount:    -10,
		Reference: "top-up",
	})

	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestDeposit_FrozenWallet(t *testing.T) {
	svc := newTestService(t)
	wallet := seedFrozenWallet(t, "owner-1", "EUR", 100)

	_, err := svc.Deposit(context.Background(), &walletv1.DepositRequest{
		WalletId:  wallet.GetID(),
		Amount:    50,
		Reference: "top-up",
	})

	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

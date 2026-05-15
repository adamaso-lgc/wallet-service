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

func TestFreezeWallet_Success(t *testing.T) {
	svc := newTestService(t)
	wallet := seedWallet(t, "owner-1", "EUR", 100)

	_, err := svc.FreezeWallet(context.Background(), &walletv1.FreezeWalletRequest{
		WalletId:  wallet.GetID(),
		Reference: "fraud-detected",
	})

	require.NoError(t, err)
	assert.Equal(t, "frozen", getView(t, svc, wallet.GetID()).Status)
}

func TestFreezeWallet_WalletNotFound(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.FreezeWallet(context.Background(), &walletv1.FreezeWalletRequest{
		WalletId:  "00000000-0000-0000-0000-000000000000",
		Reference: "fraud-detected",
	})

	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestFreezeWallet_AlreadyFrozen(t *testing.T) {
	svc := newTestService(t)
	wallet := seedFrozenWallet(t, "owner-1", "EUR", 100)

	_, err := svc.FreezeWallet(context.Background(), &walletv1.FreezeWalletRequest{
		WalletId:  wallet.GetID(),
		Reference: "fraud-detected",
	})

	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

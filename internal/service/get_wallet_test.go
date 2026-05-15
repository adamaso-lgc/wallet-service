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

func TestGetWallet_Success(t *testing.T) {
	svc := newTestService(t)
	wallet := seedWallet(t, "owner-1", "EUR", 150)

	resp, err := svc.GetWallet(context.Background(), &walletv1.GetWalletRequest{
		WalletId: wallet.GetID(),
	})

	require.NoError(t, err)
	assert.Equal(t, wallet.GetID(), resp.WalletId)
	assert.Equal(t, "owner-1", resp.OwnerId)
	assert.Equal(t, "EUR", resp.Currency)
	assert.Equal(t, float64(150), resp.Balance)
	assert.Equal(t, "active", resp.Status)
	assert.NotNil(t, resp.CreatedAt)
	assert.NotNil(t, resp.UpdatedAt)
}

func TestGetWallet_NotFound(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.GetWallet(context.Background(), &walletv1.GetWalletRequest{
		WalletId: "00000000-0000-0000-0000-000000000000",
	})

	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

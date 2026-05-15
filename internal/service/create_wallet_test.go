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

func TestCreateWallet_Success(t *testing.T) {
	svc := newTestService(t)

	resp, err := svc.CreateWallet(context.Background(), &walletv1.CreateWalletRequest{
		OwnerId:        "owner-1",
		Currency:       "EUR",
		InitialBalance: 100,
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp.WalletId)

	view := getView(t, svc, resp.WalletId)
	assert.Equal(t, "owner-1", view.OwnerId)
	assert.Equal(t, "EUR", view.Currency)
	assert.Equal(t, float64(100), view.Balance)
	assert.Equal(t, "active", view.Status)
}

func TestCreateWallet_ZeroInitialBalance(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.CreateWallet(context.Background(), &walletv1.CreateWalletRequest{
		OwnerId:  "owner-1",
		Currency: "EUR",
	})

	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateWallet_MissingOwner(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.CreateWallet(context.Background(), &walletv1.CreateWalletRequest{
		Currency:       "EUR",
		InitialBalance: 100,
	})

	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateWallet_MissingCurrency(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.CreateWallet(context.Background(), &walletv1.CreateWalletRequest{
		OwnerId:        "owner-1",
		InitialBalance: 100,
	})

	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateWallet_NegativeInitialBalance(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.CreateWallet(context.Background(), &walletv1.CreateWalletRequest{
		OwnerId:        "owner-1",
		Currency:       "EUR",
		InitialBalance: -50,
	})

	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

package postgres

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adamaso/wallet-service/internal/domain"
)

// --- buildInsertParams ---

func TestBuildInsertParams_SingleWalletSingleEvent(t *testing.T) {
	repo := NewWalletRepository(nil)
	wallet := mustWallet(t, "owner-1", "USD", 100)
	events := wallet.GetUncommittedEvents()

	params, err := repo.buildInsertParams([]walletItems{{wallet: wallet, events: events}})

	require.NoError(t, err)
	assert.Len(t, params.AggregateIds, 1)
	assert.Len(t, params.EventTypes, 1)
	assert.Len(t, params.Payloads, 1)
	assert.Len(t, params.Versions, 1)
	assert.Len(t, params.OccurredAts, 1)
	assert.True(t, params.AggregateIds[0].Valid)
	assert.Equal(t, string(domain.EventWalletCreated), params.EventTypes[0])
	assert.Equal(t, int64(1), params.Versions[0])
}

func TestBuildInsertParams_VersionSequence(t *testing.T) {
	repo := NewWalletRepository(nil)
	wallet := mustWallet(t, "owner-1", "USD", 100)
	require.NoError(t, wallet.Deposit(50, "top-up"))
	// wallet now has 2 uncommitted events (Created + Deposited), version=2
	events := wallet.GetUncommittedEvents()

	params, err := repo.buildInsertParams([]walletItems{{wallet: wallet, events: events}})

	require.NoError(t, err)
	require.Len(t, params.Versions, 2)
	assert.Equal(t, int64(1), params.Versions[0])
	assert.Equal(t, int64(2), params.Versions[1])
}

func TestBuildInsertParams_MultipleWallets(t *testing.T) {
	repo := NewWalletRepository(nil)
	w1 := mustWallet(t, "owner-1", "USD", 100) // 1 event
	w2 := mustWallet(t, "owner-2", "EUR", 200) // 1 event

	items := []walletItems{
		{wallet: w1, events: w1.GetUncommittedEvents()},
		{wallet: w2, events: w2.GetUncommittedEvents()},
	}

	params, err := repo.buildInsertParams(items)

	require.NoError(t, err)
	assert.Len(t, params.AggregateIds, 2, "one row per event across all wallets")
	assert.Len(t, params.Versions, 2)
	assert.Equal(t, int64(1), params.Versions[0])
	assert.Equal(t, int64(1), params.Versions[1])
}

func TestBuildInsertParams_VersionContinuesFromExistingHistory(t *testing.T) {
	repo := NewWalletRepository(nil)
	wallet := mustWallet(t, "owner-1", "USD", 100)
	// Simulate the wallet having been saved once: version=1, no uncommitted events.
	wallet.ClearUncommittedEvents()

	require.NoError(t, wallet.Deposit(50, "top-up"))
	// version=2, one new uncommitted event → startVersion = 2-1 = 1, so version slot = 2
	events := wallet.GetUncommittedEvents()

	params, err := repo.buildInsertParams([]walletItems{{wallet: wallet, events: events}})

	require.NoError(t, err)
	require.Len(t, params.Versions, 1)
	assert.Equal(t, int64(2), params.Versions[0], "version must continue from saved history, not restart at 1")
}

func TestMapDBError_UniqueViolation(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23505"}

	result := mapDBError(pgErr)

	require.ErrorIs(t, result, domain.ErrConcurrentModification)
}

func TestMapDBError_OtherPostgresError(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "42000", Message: "syntax error"}

	result := mapDBError(pgErr)

	// Other PG errors must pass through unchanged.
	var target *pgconn.PgError
	assert.True(t, errors.As(result, &target))
	assert.Equal(t, "42000", target.Code)
}

func TestMapDBError_NonPostgresError(t *testing.T) {
	original := fmt.Errorf("some other error")

	result := mapDBError(original)

	assert.Equal(t, original, result)
}

func TestMapDBError_Nil(t *testing.T) {
	assert.NoError(t, mapDBError(nil))
}

// --- helpers ---

func mustWallet(t *testing.T, ownerID, currency string, balance float64) *domain.Wallet {
	t.Helper()
	w, err := domain.NewWallet(ownerID, currency, balance)
	require.NoError(t, err)
	return w
}

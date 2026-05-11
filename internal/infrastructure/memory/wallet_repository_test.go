package memory_test

import (
	"context"
	"sync"
	"testing"

	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalletRepository_SaveAndGet_RoundTrip(t *testing.T) {
	repo, bus := newTestRepo()
	defer bus.Close()

	wallet := mustWallet(t)
	require.NoError(t, repo.Save(context.Background(), wallet))

	retrieved, err := repo.Get(context.Background(), wallet.GetID())
	require.NoError(t, err)
	assert.Equal(t, wallet.GetID(), retrieved.GetID())
}

func TestWalletRepository_Get_NotFound(t *testing.T) {
	repo, bus := newTestRepo()
	defer bus.Close()

	_, err := repo.Get(context.Background(), "nonexistent-id")

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestWalletRepository_Save_PublishesEvents(t *testing.T) {
	store := memory.NewEventStore()
	bus := memory.NewEventBus(10)
	repo := memory.NewWalletRepository(store, bus)

	var mu sync.Mutex
	var received []domain.Event

	bus.Subscribe(func(_ context.Context, e domain.Event) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, e)
	})

	wallet := mustWallet(t)
	eventCount := len(wallet.GetUncommittedEvents())

	require.NoError(t, repo.Save(context.Background(), wallet))
	bus.Close()

	mu.Lock()
	defer mu.Unlock()
	assert.Len(t, received, eventCount, "all uncommitted events should be published to the bus")
}

func TestWalletRepository_Save_ClearsUncommittedEvents(t *testing.T) {
	repo, bus := newTestRepo()
	defer bus.Close()

	wallet := mustWallet(t)
	require.NoError(t, repo.Save(context.Background(), wallet))

	assert.Empty(t, wallet.GetUncommittedEvents())
}

func TestWalletRepository_Save_Noop_WhenNoUncommittedEvents(t *testing.T) {
	store := memory.NewEventStore()
	bus := memory.NewEventBus(10)
	repo := memory.NewWalletRepository(store, bus)

	wallet := mustWallet(t)
	wallet.ClearUncommittedEvents() // nothing to commit

	var handlerCalled bool
	bus.Subscribe(func(_ context.Context, _ domain.Event) {
		handlerCalled = true
	})

	require.NoError(t, repo.Save(context.Background(), wallet))
	bus.Close()

	assert.False(t, handlerCalled, "no events should be published when there is nothing to save")
}

// newTestRepo creates a WalletRepository wired to fresh store and bus.
func newTestRepo() (*memory.WalletRepository, *memory.EventBus) {
	store := memory.NewEventStore()
	bus := memory.NewEventBus(10)
	return memory.NewWalletRepository(store, bus), bus
}

// mustWallet is a test helper that creates a Wallet or fails immediately.
func mustWallet(t *testing.T) *domain.Wallet {
	t.Helper()
	w, err := domain.NewWallet("owner-1", "USD", 100)
	require.NoError(t, err)
	return w
}

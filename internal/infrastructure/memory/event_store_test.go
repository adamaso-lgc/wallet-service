package memory_test

import (
	"sync"
	"testing"

	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventStore_AppendAndLoad(t *testing.T) {
	store := memory.NewEventStore()
	wallet := mustWallet(t)
	events := wallet.GetUncommittedEvents()

	require.NoError(t, store.Append(wallet.GetID(), events))

	loaded, err := store.Load(wallet.GetID())
	require.NoError(t, err)
	assert.Len(t, loaded, len(events))
}

func TestEventStore_Load_NotFound(t *testing.T) {
	store := memory.NewEventStore()

	_, err := store.Load("nonexistent-id")

	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestEventStore_Load_ReturnsCopy(t *testing.T) {
	store := memory.NewEventStore()
	wallet := mustWallet(t)
	events := wallet.GetUncommittedEvents()
	require.NoError(t, store.Append(wallet.GetID(), events))

	first, err := store.Load(wallet.GetID())
	require.NoError(t, err)

	// Appending to the returned slice must not affect the store.
	first = append(first, first[0])

	second, err := store.Load(wallet.GetID())
	require.NoError(t, err)
	assert.Len(t, second, len(events), "external mutation of returned slice must not affect the store")
}

// TestEventStore_ConcurrentReadWrite verifies there are no data races under
// concurrent access. Run with: go test -race ./...
func TestEventStore_ConcurrentReadWrite(t *testing.T) {
	store := memory.NewEventStore()

	const goroutines = 20
	var wg sync.WaitGroup

	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := mustWallet(t)
			_ = store.Append(w.GetID(), w.GetUncommittedEvents())
			_, _ = store.Load(w.GetID())
		}()
	}

	wg.Wait()
}

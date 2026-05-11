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

func TestEventBus_HandlerReceivesEvents(t *testing.T) {
	bus := memory.NewEventBus(10)

	var mu sync.Mutex
	var received []domain.Event

	bus.Subscribe(func(_ context.Context, e domain.Event) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, e)
	})

	wallet := mustWallet(t)
	events := wallet.GetUncommittedEvents()
	bus.Publish(events)
	bus.Close() // blocks until all in-flight events are fully processed

	mu.Lock()
	defer mu.Unlock()
	assert.Len(t, received, len(events))
}

func TestEventBus_MultipleHandlers_AllReceive(t *testing.T) {
	const handlerCount = 3
	bus := memory.NewEventBus(10)

	counts := make([]int, handlerCount)
	var mu sync.Mutex

	for i := range handlerCount {
		idx := i
		bus.Subscribe(func(_ context.Context, _ domain.Event) {
			mu.Lock()
			defer mu.Unlock()
			counts[idx]++
		})
	}

	wallet := mustWallet(t)
	bus.Publish(wallet.GetUncommittedEvents())
	bus.Close()

	mu.Lock()
	defer mu.Unlock()
	for i, count := range counts {
		assert.Equal(t, 1, count, "handler %d should have received exactly one event", i)
	}
}

func TestEventBus_Close_DrainsInFlight(t *testing.T) {
	bus := memory.NewEventBus(10)

	var called bool
	var mu sync.Mutex

	bus.Subscribe(func(_ context.Context, _ domain.Event) {
		mu.Lock()
		defer mu.Unlock()
		called = true
	})

	bus.Publish(mustWallet(t).GetUncommittedEvents())
	bus.Close() // must not return until the handler above has been called

	mu.Lock()
	defer mu.Unlock()
	require.True(t, called, "handler must be called before Close returns")
}

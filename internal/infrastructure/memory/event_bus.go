package memory

import (
	"context"
	"sync"

	"github.com/adamaso/wallet-service/internal/domain"
)

// HandlerFunc is a function that processes a single domain event.
type HandlerFunc func(ctx context.Context, event domain.Event)

type EventBus struct {
	mu       sync.RWMutex
	handlers []HandlerFunc
	events   chan []domain.Event
	wg       sync.WaitGroup // tracks the run goroutine
}

func NewEventBus(bufferSize int) *EventBus {
	b := &EventBus{
		events: make(chan []domain.Event, bufferSize),
	}
	b.wg.Add(1)
	go b.run()
	return b
}

// Subscribe registers a handler. Must be called before the first Publish.
func (b *EventBus) Subscribe(h HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = append(b.handlers, h)
}

// Publish sends events to the bus. Non-blocking while buffer has space.
func (b *EventBus) Publish(events []domain.Event) {
	b.events <- events
}

// Close signals the bus to stop and waits for in-flight events to drain.
func (b *EventBus) Close() {
	close(b.events)
	b.wg.Wait()
}

// run is the background goroutine that drains the channel.
func (b *EventBus) run() {
	defer b.wg.Done()
	for batch := range b.events {
		b.dispatch(batch)
	}
}

// dispatch fans out a batch of events to all handlers concurrently,
// then blocks until every handler goroutine has returned.
func (b *EventBus) dispatch(events []domain.Event) {
	b.mu.RLock()
	handlers := make([]HandlerFunc, len(b.handlers))
	copy(handlers, b.handlers)
	b.mu.RUnlock()

	var wg sync.WaitGroup
	for _, e := range events {
		for _, h := range handlers {
			wg.Add(1)
			go func(handler HandlerFunc, event domain.Event) {
				defer wg.Done()
				handler(context.Background(), event)
			}(h, e)
		}
	}
	wg.Wait()
}

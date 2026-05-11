package memory

import (
	"fmt"
	"sync"

	"github.com/adamaso/wallet-service/internal/domain"
)

// EventStore is a thread-safe, in-memory event store.
// sync.RWMutex allows multiple concurrent readers but only one writer at a time.
type EventStore struct {
	mu     sync.RWMutex
	events map[string][]domain.Event
}

func NewEventStore() *EventStore {
	return &EventStore{
		events: make(map[string][]domain.Event),
	}
}

// Append persists a batch of events for an aggregate.
func (s *EventStore) Append(aggregateID string, events []domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[aggregateID] = append(s.events[aggregateID], events...)
	return nil
}

// Load retrieves all events for an aggregate.
// Returns a copy to prevent external mutation of internal state.
func (s *EventStore) Load(aggregateID string) ([]domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, ok := s.events[aggregateID]
	if !ok || len(events) == 0 {
		return nil, fmt.Errorf("%w: %s", domain.ErrNotFound, aggregateID)
	}

	result := make([]domain.Event, len(events))
	copy(result, events)
	return result, nil
}

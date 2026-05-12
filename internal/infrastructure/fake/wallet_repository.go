package fake

import (
	"context"
	"fmt"
	"sync"

	"github.com/adamaso/wallet-service/internal/domain"
)

// WalletRepository is a thread-safe in-memory implementation of
// domain.WalletRepository for use in unit tests. It stores events per
// aggregate and replays them via domain.NewWalletFromHistory.
type WalletRepository struct {
	mu     sync.RWMutex
	events map[string][]domain.Event
}

// Compile-time check that WalletRepository satisfies domain.WalletRepository.
var _ domain.WalletRepository = (*WalletRepository)(nil)

func NewWalletRepository() *WalletRepository {
	return &WalletRepository{events: make(map[string][]domain.Event)}
}

func (r *WalletRepository) Save(ctx context.Context, w *domain.Wallet) error {
	return r.SaveAll(ctx, w)
}

func (r *WalletRepository) SaveAll(_ context.Context, wallets ...*domain.Wallet) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, w := range wallets {
		events := w.GetUncommittedEvents()
		r.events[w.GetID()] = append(r.events[w.GetID()], events...)
		w.ClearUncommittedEvents()
	}
	return nil
}

func (r *WalletRepository) Get(_ context.Context, id string) (*domain.Wallet, error) {
	r.mu.RLock()
	events, ok := r.events[id]
	r.mu.RUnlock()
	if !ok || len(events) == 0 {
		return nil, fmt.Errorf("%w: %s", domain.ErrNotFound, id)
	}
	return domain.NewWalletFromHistory(events)
}

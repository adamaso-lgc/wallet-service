package memory

import (
	"context"

	"github.com/adamaso/wallet-service/internal/domain"
)

// WalletRepository implements domain.WalletRepository using the in-memory EventStore and EventBus.
type WalletRepository struct {
	store *EventStore
	bus   *EventBus
}

// Compile-time check that WalletRepository satisfies domain.WalletRepository.
var _ domain.WalletRepository = (*WalletRepository)(nil)

func NewWalletRepository(store *EventStore, bus *EventBus) *WalletRepository {
	return &WalletRepository{store: store, bus: bus}
}

func (r *WalletRepository) Save(ctx context.Context, wallet *domain.Wallet) error {
	events := wallet.GetUncommittedEvents()
	if len(events) == 0 {
		return nil
	}
	if err := r.store.Append(wallet.GetID(), events); err != nil {
		return err
	}
	wallet.ClearUncommittedEvents()
	r.bus.Publish(events)
	return nil
}

func (r *WalletRepository) Get(ctx context.Context, id string) (*domain.Wallet, error) {
	events, err := r.store.Load(id)
	if err != nil {
		return nil, err
	}
	return domain.NewWalletFromHistory(events)
}

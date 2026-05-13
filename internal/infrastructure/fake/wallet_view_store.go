package fake

import (
	"context"
	"fmt"

	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/projection"
)

// WalletViewStore is an in-memory implementation of projection.WalletStore
// for use in unit tests.
type WalletViewStore struct {
	wallets map[string]*projection.WalletView
}

// Compile-time check that WalletViewStore satisfies projection.WalletStore.
var _ projection.WalletStore = (*WalletViewStore)(nil)

func NewWalletViewStore(views ...*projection.WalletView) *WalletViewStore {
	s := &WalletViewStore{wallets: make(map[string]*projection.WalletView)}
	for _, v := range views {
		s.wallets[v.ID] = v
	}
	return s
}

func (s *WalletViewStore) GetWallet(_ context.Context, id string) (*projection.WalletView, error) {
	v, ok := s.wallets[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", domain.ErrNotFound, id)
	}
	return v, nil
}

func (s *WalletViewStore) ListWalletsByOwner(_ context.Context, ownerID string) ([]*projection.WalletView, error) {
	var result []*projection.WalletView
	for _, v := range s.wallets {
		if v.OwnerID == ownerID {
			result = append(result, v)
		}
	}
	return result, nil
}

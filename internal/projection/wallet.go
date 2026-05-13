package projection

import (
	"context"
	"time"
)

type WalletView struct {
	ID        string
	OwnerID   string
	Balance   float64
	Currency  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WalletStore interface {
	GetWallet(ctx context.Context, id string) (*WalletView, error)
	ListWalletsByOwner(ctx context.Context, ownerID string) ([]*WalletView, error)
}

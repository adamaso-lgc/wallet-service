package domain

import "context"

type Repository[T Aggregate] interface {
	Save(ctx context.Context, aggregate T) error
	SaveAll(ctx context.Context, aggregates ...T) error
	Get(ctx context.Context, id string) (T, error)
}

// WalletRepository is a convenience alias for the wallet-specific repository.
type WalletRepository = Repository[*Wallet]

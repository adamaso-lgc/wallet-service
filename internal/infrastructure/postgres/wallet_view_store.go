package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/postgres/db"
	"github.com/adamaso/wallet-service/internal/projection"
)

type WalletViewStore struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// Compile-time check that WalletViewStore satisfies projection.WalletStore.
var _ projection.WalletStore = (*WalletViewStore)(nil)

func NewWalletViewStore(pool *pgxpool.Pool) *WalletViewStore {
	return &WalletViewStore{pool: pool, queries: db.New()}
}

func (s *WalletViewStore) GetWallet(ctx context.Context, id string) (*projection.WalletView, error) {
	aggUUID, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("parse wallet id: %w", err)
	}

	row, err := s.queries.GetWalletView(ctx, s.pool, aggUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", domain.ErrNotFound, id)
		}
		return nil, fmt.Errorf("get wallet view: %w", err)
	}

	return toWalletView(row)
}

func (s *WalletViewStore) ListWalletsByOwner(ctx context.Context, ownerID string) ([]*projection.WalletView, error) {
	rows, err := s.queries.ListWalletViewsByOwner(ctx, s.pool, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list wallet views: %w", err)
	}

	views := make([]*projection.WalletView, 0, len(rows))
	for _, row := range rows {
		v, err := toWalletView(row)
		if err != nil {
			return nil, err
		}
		views = append(views, v)
	}
	return views, nil
}

func toWalletView(row *db.WalletView) (*projection.WalletView, error) {
	balance, err := numericToFloat64(row.Balance)
	if err != nil {
		return nil, fmt.Errorf("convert balance: %w", err)
	}

	return &projection.WalletView{
		ID:        uuidToString(row.ID),
		OwnerID:   row.OwnerID,
		Balance:   balance,
		Currency:  row.Currency,
		Status:    row.Status,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

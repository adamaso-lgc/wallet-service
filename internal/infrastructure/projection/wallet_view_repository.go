package projection

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/postgres"
	"github.com/adamaso/wallet-service/internal/infrastructure/postgres/db"
)

type Repository interface {
	GetWallet(ctx context.Context, id string) (*WalletView, error)
	ListWalletsByOwner(ctx context.Context, ownerID string) ([]*WalletView, error)
}

// Compile-time check that WalletViewRepository satisfies Repository.
var _ Repository = (*WalletViewRepository)(nil)

type WalletViewRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewWalletViewRepository(pool *pgxpool.Pool) *WalletViewRepository {
	return &WalletViewRepository{pool: pool, queries: db.New()}
}

func (s *WalletViewRepository) GetWallet(ctx context.Context, id string) (*WalletView, error) {
	aggUUID, err := postgres.UUIDFromString(id)
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

func (s *WalletViewRepository) ListWalletsByOwner(ctx context.Context, ownerID string) ([]*WalletView, error) {
	rows, err := s.queries.ListWalletViewsByOwner(ctx, s.pool, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list wallet views: %w", err)
	}

	views := make([]*WalletView, 0, len(rows))
	for _, row := range rows {
		v, err := toWalletView(row)
		if err != nil {
			return nil, err
		}
		views = append(views, v)
	}
	return views, nil
}

func toWalletView(row *db.WalletView) (*WalletView, error) {
	balance, err := postgres.NumericToFloat64(row.Balance)
	if err != nil {
		return nil, fmt.Errorf("convert balance: %w", err)
	}

	return &WalletView{
		ID:        postgres.UUIDToString(row.ID),
		OwnerID:   row.OwnerID,
		Balance:   balance,
		Currency:  row.Currency,
		Status:    row.Status,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/postgres/db"
)

type WalletProjector struct {
	queries *db.Queries
}

func NewWalletProjector() *WalletProjector {
	return &WalletProjector{queries: db.New()}
}

// Apply dispatches a single event to the appropriate projection handler.
func (p *WalletProjector) Apply(ctx context.Context, dbtx db.DBTX, event domain.Event) error {
	switch e := event.(type) {
	case domain.WalletCreatedEvent:
		return p.onWalletCreated(ctx, dbtx, e)
	case domain.MoneyDepositedEvent:
		return p.onBalanceChanged(ctx, dbtx, e.GetAggregateID(), e.BalanceAfter, e.GetOccurredAt())
	case domain.MoneyWithdrawnEvent:
		return p.onBalanceChanged(ctx, dbtx, e.GetAggregateID(), e.BalanceAfter, e.GetOccurredAt())
	case domain.MoneyTransferredEvent:
		return p.onBalanceChanged(ctx, dbtx, e.GetAggregateID(), e.BalanceAfter, e.GetOccurredAt())
	case domain.WalletFrozenEvent:
		return p.onWalletFrozen(ctx, dbtx, e)
	default:
		return fmt.Errorf("projector: unknown event type %T", event)
	}
}

func (p *WalletProjector) onWalletCreated(ctx context.Context, dbtx db.DBTX, e domain.WalletCreatedEvent) error {
	id, err := uuidFromString(e.GetAggregateID())
	if err != nil {
		return err
	}
	balance, err := numericFromFloat64(e.InitialBalance)
	if err != nil {
		return err
	}
	ts := pgtype.Timestamptz{Time: e.GetOccurredAt(), Valid: true}
	return p.queries.UpsertWalletView(ctx, dbtx, &db.UpsertWalletViewParams{
		ID:        id,
		OwnerID:   e.OwnerID,
		Balance:   balance,
		Currency:  e.Currency,
		Status:    string(domain.StatusActive),
		CreatedAt: ts,
		UpdatedAt: ts,
	})
}

func (p *WalletProjector) onBalanceChanged(ctx context.Context, dbtx db.DBTX, aggregateID string, balanceAfter float64, updatedAt time.Time) error {
	id, err := uuidFromString(aggregateID)
	if err != nil {
		return err
	}
	balance, err := numericFromFloat64(balanceAfter)
	if err != nil {
		return err
	}
	return p.queries.UpdateWalletViewBalance(ctx, dbtx, &db.UpdateWalletViewBalanceParams{
		ID:        id,
		Balance:   balance,
		UpdatedAt: pgtype.Timestamptz{Time: updatedAt, Valid: true},
	})
}

func (p *WalletProjector) onWalletFrozen(ctx context.Context, dbtx db.DBTX, e domain.WalletFrozenEvent) error {
	id, err := uuidFromString(e.GetAggregateID())
	if err != nil {
		return err
	}
	return p.queries.UpdateWalletViewStatus(ctx, dbtx, &db.UpdateWalletViewStatusParams{
		ID:        id,
		Status:    string(domain.StatusFrozen),
		UpdatedAt: pgtype.Timestamptz{Time: e.GetOccurredAt(), Valid: true},
	})
}

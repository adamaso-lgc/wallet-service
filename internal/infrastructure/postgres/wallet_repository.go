package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/postgres/db"
)

// walletItems is an internal grouping of a wallet with its uncommitted events.
type walletItems struct {
	wallet *domain.Wallet
	events []domain.Event
}

type WalletRepository struct {
	pool      *pgxpool.Pool
	queries   *db.Queries
	codec     *Codec
	projector *WalletProjector
}

// Compile-time check that WalletRepository satisfies domain.WalletRepository.
var _ domain.WalletRepository = (*WalletRepository)(nil)

func NewWalletRepository(pool *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{
		pool:      pool,
		queries:   db.New(),
		codec:     NewCodec(),
		projector: NewWalletProjector(),
	}
}

func (r *WalletRepository) Save(ctx context.Context, wallet *domain.Wallet) error {
	return r.SaveAll(ctx, wallet)
}

func (r *WalletRepository) SaveAll(ctx context.Context, wallets ...*domain.Wallet) error {
	var collected []walletItems
	for _, w := range wallets {
		if events := w.GetUncommittedEvents(); len(events) > 0 {
			collected = append(collected, walletItems{wallet: w, events: events})
		}
	}
	if len(collected) == 0 {
		return nil
	}

	params, err := r.buildInsertParams(collected)
	if err != nil {
		return err
	}

	return pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		if err := r.queries.InsertEvents(ctx, tx, params); err != nil {
			return mapDBError(err)
		}

		for _, it := range collected {
			for _, event := range it.events {
				if err := r.projector.Apply(ctx, tx, event); err != nil {
					return fmt.Errorf("project event: %w", err)
				}
			}
		}

		for _, it := range collected {
			it.wallet.ClearUncommittedEvents()
		}
		return nil
	})
}

func (r *WalletRepository) Get(ctx context.Context, id string) (*domain.Wallet, error) {
	aggUUID, err := uuidFromString(id)
	if err != nil {
		return nil, fmt.Errorf("parse aggregate id: %w", err)
	}

	rows, err := r.queries.LoadEventsByAggregateID(ctx, r.pool, aggUUID)
	if err != nil {
		return nil, fmt.Errorf("load events: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("%w: %s", domain.ErrNotFound, id)
	}

	events := make([]domain.Event, 0, len(rows))
	for _, row := range rows {
		event, err := r.codec.Decode(row.EventType, row.Payload)
		if err != nil {
			return nil, fmt.Errorf("decode event: %w", err)
		}
		events = append(events, event)
	}
	return domain.NewWalletFromHistory(events)
}

func (r *WalletRepository) buildInsertParams(collected []walletItems) (*db.InsertEventsParams, error) {
	params := &db.InsertEventsParams{}

	for _, it := range collected {
		startVersion := it.wallet.GetVersion() - int64(len(it.events))
		for i, event := range it.events {
			version := startVersion + int64(i) + 1

			aggUUID, err := uuidFromString(event.GetAggregateID())
			if err != nil {
				return nil, fmt.Errorf("parse aggregate id: %w", err)
			}
			eventType, payload, err := r.codec.Encode(event)
			if err != nil {
				return nil, fmt.Errorf("encode event: %w", err)
			}

			params.AggregateIds = append(params.AggregateIds, aggUUID)
			params.EventTypes = append(params.EventTypes, eventType)
			params.Payloads = append(params.Payloads, payload)
			params.Versions = append(params.Versions, version)
			params.OccurredAts = append(params.OccurredAts, pgtype.Timestamptz{
				Time:  event.GetOccurredAt(),
				Valid: true,
			})
		}
	}
	return params, nil
}

// mapDBError translates PostgreSQL constraint violations to domain errors.
func mapDBError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case "23505": // unique_violation — two writers raced on same (aggregate_id, version)
			return fmt.Errorf("%w", domain.ErrConcurrentModification)
		}
	}
	return err
}

package postgres_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/postgres"
)

// baseEvent builds a minimal BaseEvent for use in event constructors.
func baseEvent(aggregateID string, eventType domain.EventType) domain.BaseEvent {
	return domain.BaseEvent{
		EventID:     "evt-1",
		Type:        eventType,
		AggregateID: aggregateID,
		OccurredAt:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
}

func TestCodec_RoundTrip(t *testing.T) {
	codec := postgres.NewCodec()
	const aggID = "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name  string
		event domain.Event
	}{
		{
			name: "WalletCreated",
			event: domain.WalletCreatedEvent{
				BaseEvent:      baseEvent(aggID, domain.EventWalletCreated),
				OwnerID:        "owner-1",
				Currency:       "USD",
				InitialBalance: 100.50,
			},
		},
		{
			name: "MoneyDeposited",
			event: domain.MoneyDepositedEvent{
				BaseEvent:    baseEvent(aggID, domain.EventWalletDeposited),
				Amount:       50.0,
				BalanceAfter: 150.50,
				Reference:    "top-up",
			},
		},
		{
			name: "MoneyWithdrawn",
			event: domain.MoneyWithdrawnEvent{
				BaseEvent:    baseEvent(aggID, domain.EventMoneyWithdrawn),
				Amount:       20.0,
				BalanceAfter: 130.50,
				Reference:    "purchase",
			},
		},
		{
			name: "MoneyTransferred",
			event: domain.MoneyTransferredEvent{
				BaseEvent:      baseEvent(aggID, domain.EventMoneyTransferred),
				Amount:         30.0,
				BalanceAfter:   100.50,
				CounterpartyID: "wallet-2",
				Direction:      "debit",
				Reference:      "payment",
			},
		},
		{
			name: "WalletFrozen",
			event: domain.WalletFrozenEvent{
				BaseEvent: baseEvent(aggID, domain.EventWalletFrozen),
				Reference: "compliance",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eventType, payload, err := codec.Encode(tc.event)
			require.NoError(t, err)
			assert.NotEmpty(t, eventType)
			assert.NotEmpty(t, payload)

			decoded, err := codec.Decode(eventType, payload)
			require.NoError(t, err)

			assert.Equal(t, tc.event.GetEventID(), decoded.GetEventID())
			assert.Equal(t, tc.event.GetEventType(), decoded.GetEventType())
			assert.Equal(t, tc.event.GetAggregateID(), decoded.GetAggregateID())
			assert.WithinDuration(t, tc.event.GetOccurredAt(), decoded.GetOccurredAt(), 0)
		})
	}
}

func TestCodec_Encode_ReturnsCorrectEventType(t *testing.T) {
	codec := postgres.NewCodec()
	event := domain.WalletCreatedEvent{
		BaseEvent: baseEvent("550e8400-e29b-41d4-a716-446655440000", domain.EventWalletCreated),
	}

	eventType, _, err := codec.Encode(event)

	require.NoError(t, err)
	assert.Equal(t, string(domain.EventWalletCreated), eventType)
}

func TestCodec_Decode_UnknownEventType(t *testing.T) {
	codec := postgres.NewCodec()

	_, err := codec.Decode("UnknownEvent", []byte(`{}`))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown event type")
}

func TestCodec_Decode_MalformedPayload(t *testing.T) {
	codec := postgres.NewCodec()

	_, err := codec.Decode(string(domain.EventWalletCreated), []byte(`not-json`))

	require.Error(t, err)
}

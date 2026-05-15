package eventstore

import (
	"encoding/json"
	"fmt"

	"github.com/adamaso/wallet-service/internal/domain"
)

// Codec encodes domain events to/from the JSONB payload stored in the events table.
// The full event struct is serialized so all fields — including BaseEvent — are
// preserved and can be used to reconstruct aggregate state on Load.
type Codec struct{}

func NewCodec() *Codec { return &Codec{} }

func (c *Codec) Encode(event domain.Event) (eventType string, payload []byte, err error) {
	payload, err = json.Marshal(event)
	return string(event.GetEventType()), payload, err
}

func (c *Codec) Decode(eventType string, payload []byte) (domain.Event, error) {
	switch domain.EventType(eventType) {
	case domain.EventWalletCreated:
		var e domain.WalletCreatedEvent
		return e, json.Unmarshal(payload, &e)
	case domain.EventWalletDeposited:
		var e domain.MoneyDepositedEvent
		return e, json.Unmarshal(payload, &e)
	case domain.EventMoneyWithdrawn:
		var e domain.MoneyWithdrawnEvent
		return e, json.Unmarshal(payload, &e)
	case domain.EventMoneyTransferred:
		var e domain.MoneyTransferredEvent
		return e, json.Unmarshal(payload, &e)
	case domain.EventWalletFrozen:
		var e domain.WalletFrozenEvent
		return e, json.Unmarshal(payload, &e)
	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
}

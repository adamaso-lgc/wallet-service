package domain

import "time"

type EventType string

type Event interface {
	GetEventID() string
	GetEventType() EventType
	GetAggregateID() string
	GetOccurredAt() time.Time
}

type BaseEvent struct {
	EventID     string    `json:"event_id"`
	Type        EventType `json:"type"`
	AggregateID string    `json:"aggregate_id"`
	OccurredAt  time.Time `json:"occurred_at"`
}

func (e BaseEvent) GetEventID() string       { return e.EventID }
func (e BaseEvent) GetEventType() EventType  { return e.Type }
func (e BaseEvent) GetAggregateID() string   { return e.AggregateID }
func (e BaseEvent) GetOccurredAt() time.Time { return e.OccurredAt }

package domain

import "time"

type EventType string

type Event interface {
	GetEventID() string
	GetEventType() EventType
	GetAggregateID() string
	GetOccurredAt() time.Time
	GetVersion() int64
}

type BaseEvent struct {
	EventID     string    `json:"event_id"`
	Type        EventType `json:"type"`
	AggregateID string    `json:"aggregate_id"`
	OccurredAt  time.Time `json:"occurred_at"`
	Version     int64     `json:"version"`
}

func (e BaseEvent) GetEventID() string       { return e.EventID }
func (e BaseEvent) GetEventType() EventType  { return e.Type }
func (e BaseEvent) GetAggregateID() string   { return e.AggregateID }
func (e BaseEvent) GetOccurredAt() time.Time { return e.OccurredAt }
func (e BaseEvent) GetVersion() int64        { return e.Version }

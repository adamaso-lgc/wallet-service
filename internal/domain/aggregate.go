package domain

type Aggregate interface {
	GetID() string
	GetVersion() int64
	GetUncommittedEvents() []Event
	ClearUncommittedEvents()
	Apply(event Event) error
}

type BaseAggregate struct {
	id                string
	version           int64
	uncommittedEvents []Event
}

func (a *BaseAggregate) GetID() string { return a.id }

func (a *BaseAggregate) GetVersion() int64 { return a.version }

func (a *BaseAggregate) GetUncommittedEvents() []Event { return a.uncommittedEvents }

func (a *BaseAggregate) ClearUncommittedEvents() { a.uncommittedEvents = []Event{} }

func (a *BaseAggregate) Raise(aggregate Aggregate, event Event) error {
	a.uncommittedEvents = append(a.uncommittedEvents, event)
	if err := aggregate.Apply(event); err != nil {
		return err
	}
	a.version++
	return nil
}

func (a *BaseAggregate) LoadFromHistory(aggregate Aggregate, events []Event) error {
	for _, event := range events {
		if err := aggregate.Apply(event); err != nil {
			return err
		}
		a.version++
	}
	return nil
}

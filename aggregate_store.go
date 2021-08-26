package finkgoes

import (
	"context"
	"errors"
)

type AggregateStore interface {
	Load(ctx context.Context, aggregateID string) (AggregateRoot, error)
	Store(ctx context.Context, aggregate AggregateRoot) error
}

type AggregateStoreBase struct {
	aggregateType AggregateType
	eventStore EventStore
}

func NewAggregateStore(eventStore EventStore, aggregateType AggregateType) (*AggregateStoreBase, error) {
	d := &AggregateStoreBase{
		aggregateType: aggregateType,
		eventStore: eventStore,
	}
	return d, nil
}

func (as *AggregateStoreBase) Load(ctx context.Context, aggregateID string) (AggregateRoot, error) {
	stream := as.aggregateType.StreamNameFor(aggregateID)
	as.eventStore.ReadEvents(stream, -1, -1)
	return nil, nil
}

func (as *AggregateStoreBase) Store(ctx context.Context, aggregate AggregateRoot) error {
	stream := as.aggregateType.StreamNameFor(aggregate.AggregateID())

	events := aggregate.GetEvents()
	if len(events) == 0 {
		return nil
	}

	if err := as.eventStore.AppendEvents(stream, -1, events); err != nil {
		return err
	}

	aggregate.ClearEvents()

	if err := as.applyEvents(aggregate, events); err != nil {
		return err
	}

	return nil
}
func (as *AggregateStoreBase) applyEvents(a AggregateRoot, events []Event) error {
	for _, event := range events {
		if err := a.ApplyEvent(event); err != nil {
			return errors.New("could not apply event")
		}
		a.SetVersion(event.Version())
	}

	return nil
}
package finkgoes

import (
	"context"
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
	as.eventStore.ReadEvents(ctx, stream, -1, -1)
	return nil, nil
}

func (as *AggregateStoreBase) Store(ctx context.Context, aggregate AggregateRoot) error {
	stream := as.aggregateType.StreamNameFor(aggregate.AggregateID())

	changes := aggregate.GetChanges()
	if len(changes) == 0 {
		return nil
	}

	if err := as.eventStore.AppendEvents(ctx, stream, aggregate.OriginalVersion(), changes); err != nil {
		return err
	}
	return nil
}

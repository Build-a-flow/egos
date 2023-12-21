package egos

import (
	"context"
)

type AggregateStore interface {
	Load(ctx context.Context, aggregate AggregateRoot, aggregateID string) error
	Store(ctx context.Context, aggregate AggregateRoot) error
	Delete(ctx context.Context, aggregate AggregateRoot, aggregateID string) error
}

type AggregateStoreBase struct {
	aggregate  AggregateRoot
	eventStore EventStore
}

func NewAggregateStore(eventStore EventStore, aggregate AggregateRoot) (*AggregateStoreBase, error) {
	d := &AggregateStoreBase{
		aggregate:  aggregate,
		eventStore: eventStore,
	}
	return d, nil
}

func (as *AggregateStoreBase) Load(ctx context.Context, aggregate AggregateRoot, aggregateID string) error {
	stream := StreamNameFor(aggregate, aggregateID)
	events := as.eventStore.ReadEvents(ctx, stream, -1, -1)
	aggregate.Fold(aggregate, events)
	return nil
}

func (as *AggregateStoreBase) Store(ctx context.Context, aggregate AggregateRoot) error {
	stream := StreamNameFor(aggregate, aggregate.AggregateID())
	changes := aggregate.GetChanges()
	if len(changes) == 0 {
		return nil
	}

	if err := as.eventStore.AppendEvents(ctx, stream, aggregate.OriginalVersion(), changes); err != nil {
		return err
	}
	return nil
}

func (as *AggregateStoreBase) Delete(ctx context.Context, aggregateID string) error {
	stream := StreamNameFor(as.aggregate, aggregateID)
	if err := as.eventStore.DeleteStream(ctx, stream); err != nil {
		return err
	}
	return nil
}

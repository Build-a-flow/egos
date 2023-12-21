package store

import (
	"context"

	"github.com/finktek/egos"
)

type AggregateStore interface {
	Load(ctx context.Context, aggregate egos.AggregateRoot, aggregateID string) error
	Store(ctx context.Context, aggregate egos.AggregateRoot) error
}

type AggregateStoreBase struct {
	aggregate  egos.AggregateRoot
	eventStore EventStore
}

func NewAggregateStore(eventStore EventStore, aggregate egos.AggregateRoot) (*AggregateStoreBase, error) {
	d := &AggregateStoreBase{
		aggregate:  aggregate,
		eventStore: eventStore,
	}
	return d, nil
}

func (as *AggregateStoreBase) Load(ctx context.Context, aggregate egos.AggregateRoot, aggregateID string) error {
	stream := egos.StreamNameFor(aggregate, aggregateID)
	events := as.eventStore.ReadEvents(ctx, stream, -1, -1)
	aggregate.Fold(aggregate, events)
	return nil
}

func (as *AggregateStoreBase) Store(ctx context.Context, aggregate egos.AggregateRoot) error {
	stream := egos.StreamNameFor(aggregate, aggregate.AggregateID())
	changes := aggregate.GetChanges()
	if len(changes) == 0 {
		return nil
	}

	if err := as.eventStore.AppendEvents(ctx, stream, aggregate.OriginalVersion(), changes); err != nil {
		return err
	}
	return nil
}

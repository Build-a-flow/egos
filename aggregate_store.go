package egos

import (
	"context"
	"math"

	"github.com/gofrs/uuid"
)

type AggregateStore interface {
	NewId(ctx context.Context) string
	Load(ctx context.Context, aggregate AggregateRoot, aggregateID string) error
	Store(ctx context.Context, aggregate AggregateRoot) error
	Delete(ctx context.Context, aggregateID string) error
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

func (as *AggregateStoreBase) NewId(ctx context.Context) string {
	return uuid.Must(uuid.NewV4()).String()
}

func (as *AggregateStoreBase) Load(ctx context.Context, aggregate AggregateRoot, aggregateID string) error {
	stream := StreamNameFor(aggregate, aggregateID)
	events, err := as.eventStore.ReadEvents(ctx, stream, 0, math.MaxInt64)
	if err != nil {
		return err
	}
	aggregate.Fold(aggregate, events)
	return nil
}

func (as *AggregateStoreBase) Store(ctx context.Context, aggregate AggregateRoot) error {
	stream := StreamNameFor(aggregate, aggregate.AggregateID())
	changes := aggregate.GetChanges()
	if len(changes) == 0 {
		return nil
	}

	if err := as.eventStore.AppendEvents(ctx, stream, int64(aggregate.OriginalVersion()), changes); err != nil {
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

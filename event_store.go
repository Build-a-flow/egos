package egos

import (
	"context"

	"github.com/build-a-flow/egos/eventstore"
)

type EventStore interface {
	AppendEvents(ctx context.Context, streamName string, expectedVersion int, events []Event) error
	ReadEvents(ctx context.Context, streamName string, start int, limit int) []Event
	ReadEventsBackwards(ctx context.Context, streamName string, limit int) []Event
	ReadStream(ctx context.Context, streamName string, start int, callback func())
	TruncateStream(ctx context.Context, streamName string, position int, expectedVersion int)
	DeleteStream(ctx context.Context, streamName string, expectedVersion int)
}

func NewEventStoreDbClient(connectionString string) (*EventStore, error) {
	eventStore, err := eventstore.NewEventStoreDbClient(connectionString)
	return eventStore, err
}

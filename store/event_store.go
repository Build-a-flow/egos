package store

import (
	"context"

	"github.com/finktek/egos"
)

type EventStore interface {
	AppendEvents(ctx context.Context, streamName string, expectedVersion int, events []egos.Event) error
	ReadEvents(ctx context.Context, streamName string, start int, limit int) []egos.Event
	ReadEventsBackwards(ctx context.Context, streamName string, limit int) []egos.Event
	ReadStream(ctx context.Context, streamName string, start int, callback func())
	TruncateStream(ctx context.Context, streamName string, position int, expectedVersion int)
	DeleteStream(ctx context.Context, streamName string, expectedVersion int)
}

package egos

import (
	"context"
)

type EventStore interface {
	AppendEvents(ctx context.Context, streamName string, expectedVersion int64, events []Event) error
	ReadEvents(ctx context.Context, streamName string, start int64, limit int64) ([]Event, error)
	DeleteStream(ctx context.Context, streamName string) error
}

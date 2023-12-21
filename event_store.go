package egos

import (
	"context"
)

type EventStore interface {
	AppendEvents(ctx context.Context, streamName string, expectedVersion int, events []Event) error
	ReadEvents(ctx context.Context, streamName string, start int, limit int) []Event
	DeleteStream(ctx context.Context, streamName string) error
}

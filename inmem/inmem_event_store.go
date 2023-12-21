package inmem

import (
	"context"

	"github.com/finktek/egos"
)

type InMemEventStore struct {
	Streams map[string][]egos.Event
}

func NewInMemEventStore() egos.EventStore {
	return &InMemEventStore{
		Streams: make(map[string][]egos.Event),
	}
}

func (s *InMemEventStore) AppendEvents(ctx context.Context, streamName string, expectedVersion int, events []egos.Event) error {
	if _, ok := s.Streams[streamName]; !ok {
		s.Streams[streamName] = make([]egos.Event, 0)
	}
	if len(s.Streams[streamName]) != expectedVersion {
		return egos.ErrConcurrencyViolation
	}
	s.Streams[streamName] = append(s.Streams[streamName], events...)
	return nil
}

func (s *InMemEventStore) ReadEvents(ctx context.Context, streamName string, start int, limit int) []egos.Event {
	if _, ok := s.Streams[streamName]; !ok {
		return nil
	}
	if start > len(s.Streams[streamName]) {
		return nil
	}
	if start+limit > len(s.Streams[streamName]) {
		return s.Streams[streamName][start:]
	}
	return s.Streams[streamName][start : start+limit]
}

func (s *InMemEventStore) DeleteStream(ctx context.Context, streamName string) error {
	delete(s.Streams, streamName)
	return nil
}

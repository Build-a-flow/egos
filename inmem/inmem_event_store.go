package inmem

import (
	"context"
	"errors"

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

func (s *InMemEventStore) AppendEvents(ctx context.Context, streamName string, expectedVersion int64, events []egos.Event) error {
	if _, ok := s.Streams[streamName]; !ok {
		s.Streams[streamName] = make([]egos.Event, 0)
	}
	s.Streams[streamName] = append(s.Streams[streamName], events...)
	return nil
}

func (s *InMemEventStore) ReadEvents(ctx context.Context, streamName string, start int64, limit int64) ([]egos.Event, error) {
	if _, ok := s.Streams[streamName]; !ok {
		return nil, errors.New("stream not found")
	}
	if start == -1 && limit == -1 {
		return s.Streams[streamName], nil
	}
	if start > int64(len(s.Streams[streamName])) {
		return nil, errors.New("start index out of bounds")
	}
	if start+limit > int64(len(s.Streams[streamName])) {
		return s.Streams[streamName][start:], nil
	}
	return s.Streams[streamName][start : start+limit], nil
}

func (s *InMemEventStore) DeleteStream(ctx context.Context, streamName string) error {
	delete(s.Streams, streamName)
	return nil
}

package esdb

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/finktek/egos"
	"github.com/gofrs/uuid"
)

type EsdbEventStore struct {
	client *esdb.Client
}

func NewEsdbEventStore(config *esdb.Configuration) (egos.EventStore, error) {
	client, err := esdb.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &EsdbEventStore{client}, nil
}

func (s *EsdbEventStore) AppendEvents(ctx context.Context, streamName string, expectedVersion int64, events []egos.Event) error {
	_, err := s.client.AppendToStream(ctx, streamName, esdb.AppendToStreamOptions{}, eventsToProposedEvents(events)...)
	if err != nil {
		return err
	}
	return nil
}

func (s *EsdbEventStore) ReadEvents(ctx context.Context, streamName string, start int64, limit int64) ([]egos.Event, error) {
	from := esdb.StreamRevision{Value: uint64(start)}
	stream, err := s.client.ReadStream(ctx, streamName, esdb.ReadStreamOptions{From: from}, uint64(limit))
	if err != nil {
		if err, ok := esdb.FromError(err); !ok {
			if err.Code() == esdb.ErrorCodeResourceNotFound {
				return nil, errors.New("EOF")
			}
		}
		return nil, err
	} else if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return resolvedEventsToEvents(stream), nil
}

func (s *EsdbEventStore) DeleteStream(ctx context.Context, streamName string) error {
	return nil
}

func eventsToProposedEvents(events []egos.Event) []esdb.EventData {
	var proposedEvents []esdb.EventData
	for _, event := range events {
		serializedData, serializedHeaders := event.Serialize()
		msg := esdb.EventData{
			EventID:     uuid.Must(uuid.NewV4()),
			EventType:   event.EventType(),
			ContentType: esdb.ContentTypeJson,
			Data:        serializedData,
			Metadata:    serializedHeaders,
		}

		proposedEvents = append(proposedEvents, msg)
	}
	return proposedEvents
}

func resolvedEventsToEvents(readStream *esdb.ReadStream) []egos.Event {
	var events []egos.Event

	for {
		event, err := readStream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			break
		}

		eventData := egos.GetEventInstance(event.Event.EventType)
		if eventData != nil {
			json.Unmarshal(event.Event.Data, &eventData)
			msg := &egos.EventDescriptor{
				Data:    eventData,
				Headers: make(map[string]interface{}),
			}
			events = append(events, msg)
		}
	}
	return events
}

package esdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	egos "github.com/finktek/egos/core"
	"github.com/gofrs/uuid"
)

type EsdbEventStore struct {
	client *esdb.Client
}

func NewEsdbEventStore(client *esdb.Client) egos.EventStore {
	return &EsdbEventStore{client}
}

func (s *EsdbEventStore) AppendEvents(ctx context.Context, streamName string, expectedVersion int64, events []egos.Event) error {
	_, err := s.client.AppendToStream(ctx, streamName, esdb.AppendToStreamOptions{}, eventsToProposedEvents(events)...)
	if err != nil {
		return err
	}
	return nil
}

func (s *EsdbEventStore) ReadEvents(ctx context.Context, streamName string, start int64, limit int64) ([]egos.Event, error) {

	ropts := esdb.ReadStreamOptions{
		Direction: esdb.Forwards,
		From:      esdb.Revision(uint64(start)),
	}

	stream, err := s.client.ReadStream(ctx, streamName, ropts, uint64(limit))

	if err != nil {
		panic(err)
	} else if ctx.Err() != nil {
		panic(err)
	}

	defer stream.Close()

	return resolvedEventsToEvents(stream), nil
}

func (s *EsdbEventStore) DeleteStream(ctx context.Context, streamName string) error {
	return nil
}

func eventsToProposedEvents(events []egos.Event) []esdb.EventData {
	var proposedEvents []esdb.EventData
	for _, event := range events {
		serializedData, serializedMetadata := event.Serialize()
		msg := esdb.EventData{
			EventID:     uuid.Must(uuid.NewV4()),
			EventType:   event.EventType(),
			ContentType: esdb.ContentTypeJson,
			Data:        serializedData,
			Metadata:    serializedMetadata,
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

		if err, ok := esdb.FromError(err); !ok {
			if err.Code() == esdb.ErrorCodeResourceNotFound {
				fmt.Print("Stream not found")
				panic(err)
			} else if errors.Is(err, io.EOF) {
				break
			} else {
				panic(err)
			}
		}

		eventData := egos.GetEventInstance(event.Event.EventType)
		if eventData != nil {
			json.Unmarshal(event.Event.Data, &eventData)
			msg := &egos.EventDescriptor{
				Data:     eventData,
				Metadata: egos.NewMetadata(),
			}
			events = append(events, msg)
		}
	}
	return events
}

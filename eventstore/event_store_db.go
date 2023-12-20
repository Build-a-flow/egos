package eventstore

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/build-a-flow/egos"
	"github.com/gofrs/uuid"
)

type EventStore struct {
	client *esdb.Client
}

func NewEventStoreDbClient(connectionString string) (*egos.EventStore, error) {

	eventStoreDbConfig, err := esdb.ParseConnectionString(connectionString)
	if err != nil {
		log.Fatalf("Unexpected configuration error: %s", err.Error())
	}
	eventStoreDbConfig.SkipCertificateVerification = true
	eventStoreDbClient, err := esdb.NewClient(eventStoreDbConfig)
	if err != nil {
		log.Fatalf("Unexpected failure setting up test connection: %s", err.Error())
	}

	return &EventStore{
		client: eventStoreDbClient,
	}, nil
}

func (es EventStore) AppendEvents(ctx context.Context, streamName string, expectedVersion int, events []egos.Event) error {
	_, err := es.client.AppendToStream(ctx, streamName, esdb.AppendToStreamOptions{}, eventsToProposedEvents(events)...)
	if err != nil {
		log.Fatalf("Unexpected failure appending events: %s", err.Error())
	}
	return nil
}
func (es EventStore) ReadEvents(ctx context.Context, streamName string, start int, limit int) []egos.Event {
	eventStream, err := es.client.ReadStream(context.Background(), streamName, esdb.ReadStreamOptions{}, 1000)
	if err != nil {
		log.Fatalf("Unexpected failure setting up test connection: %s", err.Error())
	}
	data := resolvedEventsToEvents(eventStream)
	return data
}
func (es EventStore) ReadEventsBackwards(ctx context.Context, streamName string, limit int) []egos.Event {
	return nil
}
func (es EventStore) ReadStream(ctx context.Context, streamName string, start int, callback func()) {
	callback()
}
func (es EventStore) TruncateStream(ctx context.Context, streamName string, position int, expectedVersion int) {
}
func (es EventStore) DeleteStream(ctx context.Context, streamName string, expectedVersion int) {}

func eventsToProposedEvents(events []egos.Event) []esdb.EventData {
	var proposedEvents []esdb.EventData
	for _, event := range events {
		serializedData, serializedHeaders := event.Serialize()
		msg := esdb.EventData{
			EventID:     uuid.Must(uuid.NewV4()),
			EventType:   event.EventType(),
			ContentType: esdb.JsonContentType,
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

package eventstore

import (
	"context"
	"encoding/json"
	"log"

	"github.com/EventStore/EventStore-Client-Go/client"
	"github.com/EventStore/EventStore-Client-Go/connection"
	"github.com/EventStore/EventStore-Client-Go/direction"
	"github.com/EventStore/EventStore-Client-Go/messages"
	"github.com/EventStore/EventStore-Client-Go/stream_position"
	"github.com/EventStore/EventStore-Client-Go/streamrevision"
	"github.com/build-a-flow/egos"
	"github.com/google/uuid"
)

type EventStore struct {
	client *client.Client
}

func NewEventStoreDbClient(connectionString string) (*EventStore, error) {

	eventStoreDbConfig, err := connection.ParseConnectionString(connectionString)
	if err != nil {
		log.Fatalf("Unexpected configuration error: %s", err.Error())
	}
	eventStoreDbClient, err := client.NewClient(eventStoreDbConfig)
	if err != nil {
		log.Fatalf("Unexpected failure setting up test connection: %s", err.Error())
	}

	return &EventStore{
		client: eventStoreDbClient,
	}, nil
}

func (es EventStore) AppendEvents(ctx context.Context, streamName string, expectedVersion int, events []egos.Event) error {
	_, err := es.client.AppendToStream(ctx, streamName, streamrevision.StreamRevisionAny, eventsToProposedEvents(events))
	if err != nil {
		log.Fatalf("Unexpected failure appending events: %s", err.Error())
	}
	return nil
}
func (es EventStore) ReadEvents(ctx context.Context, streamName string, start int, limit int) []egos.Event {
	events, err := es.client.ReadStreamEvents(context.Background(), direction.Forwards, streamName, stream_position.Start{}, 1000, true)
	if err != nil {
		log.Fatalf("Unexpected failure setting up test connection: %s", err.Error())
	}
	data := resolvedEventsToEvents(events)
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

func eventsToProposedEvents(events []egos.Event) []messages.ProposedEvent {
	var proposedEvents []messages.ProposedEvent
	for _, event := range events {
		serializedData, serializedHeaders := event.Serialize()
		msg := messages.ProposedEvent{
			EventID:      uuid.Must(uuid.NewV4()),
			EventType:    event.EventType(),
			ContentType:  "application/json",
			Data:         serializedData,
			UserMetadata: serializedHeaders,
		}

		proposedEvents = append(proposedEvents, msg)
	}
	return proposedEvents
}
func resolvedEventsToEvents(resolvedEvents []messages.ResolvedEvent) []egos.Event {
	var events []egos.Event
	for _, resolvedEvent := range resolvedEvents {
		eventData := egos.GetEventInstance(resolvedEvent.Event.EventType)
		if eventData != nil {
			json.Unmarshal(resolvedEvent.Event.Data, &eventData)
			msg := &egos.EventDescriptor{
				Data:    eventData,
				Headers: make(map[string]interface{}),
			}
			events = append(events, msg)
		}
	}
	return events
}

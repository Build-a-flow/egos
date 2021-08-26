package eventstore

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/connection"
	"github.com/EventStore/EventStore-Client-Go/direction"
	"github.com/EventStore/EventStore-Client-Go/messages"
	"github.com/EventStore/EventStore-Client-Go/stream_position"
	"github.com/EventStore/EventStore-Client-Go/streamrevision"
	"github.com/finktek/eventum"
	"github.com/EventStore/EventStore-Client-Go/client"
	"github.com/gofrs/uuid"
	"log"
)

type EventStore struct {
	client	*client.Client
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

func (es EventStore) AppendEvents(streamName string, expectedVersion int, events []finkgoes.Event) error {
	id, _ :=  uuid.NewV4()
	_, err := es.client.AppendToStream(context.Background(), streamName, streamrevision.StreamRevisionAny, []messages.ProposedEvent{{EventID: id, EventType: "EventTypePVZ" }})
	if err != nil {
		log.Fatalf("Unexpected failure setting up test connection: %s", err.Error())
	}
	return nil
}
func (es EventStore) ReadEvents(streamName string, start int, limit int) []finkgoes.Event {
	events, err := es.client.ReadAllEvents(context.Background(), direction.Forwards, stream_position.Start{}, 1000, true)
	if err != nil {
		log.Fatalf("Unexpected failure setting up test connection: %s", err.Error())
	}
	for i := 0; i < len(events); i++ {
		log.Println(events[i].GetOriginalEvent().EventType)
	}
	return nil
}
func (es EventStore) ReadEventsBackwards(streamName string, limit int) []finkgoes.Event { return nil }
func (es EventStore) ReadStream(streamName string, start int, callback func()) {
	callback()
}
func (es EventStore) TruncateStream(streamName string, position int, expectedVersion int) {}
func (es EventStore) DeleteStream(streamName string, expectedVersion int) {}
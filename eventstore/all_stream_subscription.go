package eventstore

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/client"
	"github.com/EventStore/EventStore-Client-Go/connection"
	"github.com/EventStore/EventStore-Client-Go/messages"
	"github.com/EventStore/EventStore-Client-Go/position"
	"github.com/EventStore/EventStore-Client-Go/stream_position"
	"github.com/finktek/eventum/subscriptions"
	"log"
)

type AllStreamSubscription struct {
	client				*client.Client
	subscriptionID		string
	checkpointStore		subscriptions.CheckpointStore
	handlers 			[]subscriptions.EventHandler
	subscription 		*client.Subscription
}

func NewAllStreamSubscription(connectionString string, subscriptionID string, checkpointStore subscriptions.CheckpointStore) (*AllStreamSubscription, error) {
	eventStoreDbConfig, err := connection.ParseConnectionString(connectionString)
	if err != nil {
		log.Fatalf("Unexpected configuration error: %s", err.Error())
	}

	eventStoreDbClient, err := client.NewClient(eventStoreDbConfig)
	if err != nil {
		log.Fatalf("Unexpected failure setting up test connection: %s", err.Error())
	}

	return &AllStreamSubscription{
		client: eventStoreDbClient,
		subscriptionID: subscriptionID,
		checkpointStore: checkpointStore,
	}, nil
}

func (b *AllStreamSubscription) Start(ctx context.Context) {
	checkpoint := b.checkpointStore.GetLastCheckpoint(b.subscriptionID)
	b.subscription = b.subscribe(ctx, checkpoint)
}

func (b *AllStreamSubscription) Stop() {
	if err := b.subscription.Close(); err != nil {
		log.Fatalf("error closing subscription: %s", err)
	}
}

func (b *AllStreamSubscription) AddHandler(handler subscriptions.EventHandler) {
	b.handlers = append(b.handlers, handler)
}

func (b *AllStreamSubscription) subscribe(ctx context.Context, checkpoint subscriptions.Checkpoint) *client.Subscription {
	streamPos := position.Position{Commit: checkpoint.Position, Prepare: checkpoint.Position}
	subscription, err := b.client.SubscribeToAll(ctx, stream_position.Position{Value: streamPos}, false)
	if err != nil {
		log.Fatalf("error subscribing %s", err)
	}
	go func() {
		for {
			subEvent := subscription.Recv()

			if subEvent.Dropped != nil {
				break
			}

			if subEvent.EventAppeared != nil {
				event := subEvent.EventAppeared
				for _, h := range b.handlers {
					if h != nil {
						log.Println(event.GetOriginalEvent().Position)
						data := resolvedEventsToEvents([]messages.ResolvedEvent{*event})
						if len(data) > 0 {
							h.Handle(data[0])
						}
					}
				}
			}
		}
	}()

	return subscription
}
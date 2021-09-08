package eventstore

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/client"
	"github.com/EventStore/EventStore-Client-Go/connection"
	"github.com/EventStore/EventStore-Client-Go/messages"
	"github.com/EventStore/EventStore-Client-Go/stream_position"
	"github.com/finktek/eventum/subscriptions"
	"log"
)

type AllStreamSubscription struct {
	client	*client.Client
	handlers []subscriptions.EventHandler
	subscription *client.Subscription
}

func NewAllStreamSubscription(connectionString string) (*AllStreamSubscription, error) {
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
	}, nil
}

func (b *AllStreamSubscription) Start(ctx context.Context) {
	b.subscription = b.subscribe(ctx)
}

func (b *AllStreamSubscription) Stop() {
	if err := b.subscription.Close(); err != nil {
		log.Fatalf("error closing subscription: %s", err)
	}
}

func (b *AllStreamSubscription) AddHandler(handler subscriptions.EventHandler) {
	b.handlers = append(b.handlers, handler)
}

func (b *AllStreamSubscription) subscribe(ctx context.Context) *client.Subscription {
	subscription, _ := b.client.SubscribeToAll(ctx, stream_position.Start{}, false)

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
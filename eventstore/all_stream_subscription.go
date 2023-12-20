package eventstore

import (
	"context"
	"log"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/build-a-flow/egos/subscriptions"
)

type AllStreamSubscription struct {
	client          *esdb.Client
	subscriptionID  string
	checkpointStore subscriptions.CheckpointStore
	handlers        []subscriptions.EventHandler
	subscription    *esdb.Subscription
}

func NewAllStreamSubscription(connectionString string, subscriptionID string, checkpointStore subscriptions.CheckpointStore) (*AllStreamSubscription, error) {
	eventStoreDbConfig, err := esdb.ParseConnectionString(connectionString)
	if err != nil {
		log.Fatalf("Unexpected configuration error: %s", err.Error())
	}

	eventStoreDbClient, err := esdb.NewClient(eventStoreDbConfig)
	if err != nil {
		log.Fatalf("Unexpected failure setting up test connection: %s", err.Error())
	}

	return &AllStreamSubscription{
		client:          eventStoreDbClient,
		subscriptionID:  subscriptionID,
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

func (b *AllStreamSubscription) subscribe(ctx context.Context, checkpoint subscriptions.Checkpoint) *esdb.Subscription {
	streamPos := esdb.Position{Commit: checkpoint.Position, Prepare: checkpoint.Position}
	subscription, err := b.client.SubscribeToAll(ctx, esdb.SubscribeToAllOptions{From: streamPos})
	if err != nil {
		log.Fatalf("error subscribing %s", err)
	}
	go func() {
		for {
			subEvent := subscription.Recv()

			if subEvent.SubscriptionDropped != nil {
				break
			}

			if subEvent.EventAppeared != nil {
				event := subEvent.EventAppeared
				for _, h := range b.handlers {
					if h != nil {
						log.Println(event.Event.Position)
						//data := resolvedEventsToEvents([]esdb.ResolvedEvent{*event})
						//if len(data) > 0 {
						//		h.Handle(data[0])
						//	}
					}
				}
			}
		}
	}()

	return subscription
}
